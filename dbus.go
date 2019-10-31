package main

import (
	"context"
	"fmt"
	"os"

	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/v5"
)

const (
	UrgencyLow      Urgency = 0x0
	UrgencyNormal   Urgency = 0x1
	UrgencyCritical Urgency = 0x2

	DBusContract string = `
	<node name="org/freedesktop/Notifications">
		<interface name="org.freedesktop.Notifications">
			<method name="GetCapabilities">
			        <arg direction="out" name="capabilities"    type="as"/>
			</method>

			<method name="Notify">
				<arg direction="in" name="app_name" type="s"/>
				<arg direction="in" name="replaces_id" type="u"/>
				<arg direction="in" name="app_icon" type="s"/>
				<arg direction="in" name="summary" type="s"/>
				<arg direction="in" name="body" type="s"/>
				<arg direction="in" name="actions" type="as"/>
				<arg direction="in" name="hints" type="a{sv}"/>
				<arg direction="in" name="expire_timeout" type="i"/>
				<arg direction="out" name="id" type="u"/>
			</method>
		   	<method name="CloseNotification">
		                <arg direction="in"  name="id"              type="u"/>
			</method>"

			<method name="GetServerInformation">
		       		<arg direction="out" name="name"            type="s"/>
				<arg direction="out" name="vendor"          type="s"/>
			        <arg direction="out" name="version"         type="s"/>
			        <arg direction="out" name="spec_version"    type="s"/>
			</method>

			<signal name="NotificationClosed">
				<arg name="id"         type="u"/>
			        <arg name="reason"     type="u"/>
		        </signal>

	                <signal name="ActionInvoked">
		 	       <arg name="id"         type="u"/>
		               <arg name="action_key" type="s"/>
		        </signal>
		</interface>
	</node>
	`
)

type (
	MonitorMessage struct {
		Summary string
		Body    string
		Urgency Urgency
	}
	Urgency     uint8
	DBusHandler struct {
		queue chan<- MonitorMessage
	}
)

func (dh *DBusHandler) Notify(
	app_name string,
	replaces_id uint32,
	app_icon string,
	summary string,
	body string,
	actions []string,
	hints map[string]interface{},
	expire_timeout int32,
) (uint, *dbus.Error) {

	m := MonitorMessage{
		Summary: summary,
		Body: body,
		Urgency: Urgency(hints["urgency"].(uint8)),
	}
	dh.queue <- m
	return 1, nil
}

func (dh *DBusHandler) GetCapabilities() ([]string, *dbus.Error) {
	return []string{
		"body",
	}, nil
}

func (dh *DBusHandler) CloseNotification(id uint32) *dbus.Error {
	return nil
}

func (dh *DBusHandler) GetServerInformation() (string, string, string, string, *dbus.Error) {
	name := "dwmstatus"
	vendor := "Andreas Malmqvist"
	version := "v1.1"
	spec_version := "1.2"
	return name, vendor, version, spec_version, nil
}

// func (dh *DBusHandler) GetServerInformation(name *string, vendor *string, version *string, spec_version *string) *dbus.Error {
// 	*name = "dwmstatus"
// 	*vendor = "Andreas Malmqvist"
// 	*version = "v1.1"
// 	*spec_version = "1.2"
// 	return nil
// }

func (dh *DBusHandler) Handle(ctx context.Context, mc chan<- MonitorMessage) {
	// Connect to dbus
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not connect to dbus: %v\n", err)
		return
	}
	defer conn.Close()

	// Set monitor message channel
	dh.queue = mc

	// Export methods
	conn.Export(dh, "/org/freedesktop/Notifications", "org.freedesktop.Notifications")
	conn.Export(introspect.Introspectable(DBusContract), "/org/freedesktop/Notifications", "org.freedesktop.DBus.Introspectable")

	// Register name on dbus
	reply, err := conn.RequestName("org.freedesktop.Notifications", dbus.NameFlagDoNotQueue)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not request name: %v\n", err)
		return
	}

	// Check to see if notification system already exist
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintf(os.Stderr, "Notification system already up, skipping\n")
		return
	}

	fmt.Println("Listening for notifications..")

	// Wait until context is done
	<-ctx.Done()

	fmt.Println("Shutting down notifications daemon")
}

func (m MonitorMessage) String() string {
	return fmt.Sprintf("[%s] %s: %s", m.Urgency, m.Summary, m.Body)
}

func (u Urgency) String() string {
	v := []string{
		"LOW",
		"NORMAL",
		"CRITICAL",
	}
	return v[u]
}

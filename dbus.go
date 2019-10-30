package main

import (
	"context"
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	UrgencyLow      Urgency = 0x0
	UrgencyNormal   Urgency = 0x1
	UrgencyCritical Urgency = 0x2
)

type (
	MonitorMessage struct {
		Summary string
		Body    string
		Urgency Urgency
	}
	Urgency uint8
)

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

func NotifyMonitor(ctx context.Context, mc chan<- MonitorMessage) error {
	var (
		rules = []string{
			"type='signal',member='Notify',path='/org/freedesktop/Notifications',interface='org.freedesktop.Notifications'",
			"type='method_call',member='Notify',path='/org/freedesktop/Notifications',interface='org.freedesktop.Notifications'",
			"type='method_return',member='Notify',path='/org/freedesktop/Notifications',interface='org.freedesktop.Notifications'",
			"type='error',member='Notify',path='/org/freedesktop/Notifications',interface='org.freedesktop.Notifications'",
		}
		flag = uint(0)
	)

	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	call := conn.BusObject().Call("org.freedesktop.DBus.Monitoring.BecomeMonitor", 0, rules, flag)
	if call.Err != nil {
		return err
	}

	c := make(chan *dbus.Message, 10)
	conn.Eavesdrop(c)

	for {
		select {
		case v := <-c:
			if v.Type != dbus.TypeMethodCall {
				continue
			}
			opts := v.Body[6].(map[string]dbus.Variant)
			summary := v.Body[3].(string)
			body := v.Body[4].(string)
			urgency := opts["urgency"].Value().(uint8)
			m := MonitorMessage{
				Summary: summary,
				Body:    body,
				Urgency: Urgency(urgency),
			}
			mc <- m
		case <-ctx.Done():
			return nil
		}
	}
}

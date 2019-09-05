package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Event struct {
	Name  string
	Value string
	Order int
}

type Events []Event

func (e Events) Len() int { return len(e) }
func (e Events) Less(i, j int) bool {
	return e[i].Order < e[j].Order
}
func (e Events) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func SendUpdate(target string) {
	client, err := rpc.DialHTTP("tcp", "localhost:8910")
	if err != nil {
		panic(err)
	}
	status := new(bool) 
	err = client.Call("Server.Call", Args{Target: target}, status)
	if err != nil {
		panic(err)
	}
	if !*status {
		fmt.Printf("Request failed ...\n")
	}
}

func main() {
	// Check if user wants to update daemon
	updateFlag := flag.String("update", "", "update item with name [NAME]")
	flag.Parse()
	if len(*updateFlag) != 0 {
		SendUpdate(*updateFlag)
		return
	}
	// Load config
	config, err := NewConfig(os.Getenv("HOME") + "/.config/dwmstatus/config.yaml")
	if err != nil {
		panic(err)
	}
	// Create item cache
	cache := make(map[string]Event)
	// Create buffered event channel
	c := make(chan Event)
	// Create map for holding item specific channels
	remotes := make(map[string]chan bool)
	// Initiate workers for each item
	for order, item := range config.Items {
		remote := make(chan bool)
		remotes[item.Name] = remote
		go Worker(item, c, order, remote)
	}
	// Setup rpc server
	server := NewServer(remotes)
	rpc.Register(server)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":8910")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	go http.Serve(l, nil)
	// Listen for events and update root window
	for {
		// Create event slice to for sorting
		events := make(Events, 0)
		// Get next event
		e := <-c
		events = append(events, e)
		// Add event to cache
		cache[e.Name] = e
		// Gather rest from cache
		for _, ee := range cache {
			if ee.Name != e.Name {
				events = append(events, ee)
			}
		}
		// Sort list of events
		sort.Sort(events)
		// Store values in buffer
		buf := new(bytes.Buffer)
		for _, ee := range events {
			fmt.Fprintf(buf, "%s ", strings.TrimSpace(ee.Value))
		}
		// Update root window
		cmd := exec.Command("xsetroot", "-name", buf.String())
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
}

func Worker(item Item, c chan<- Event, order int, remote <-chan bool) {
	ticker := time.NewTicker(time.Second * item.Interval)
	for {
		// Run command
		b, err := RunCommand(item.Command)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not issue command %s: %v\n", item.Name, err)
			return
		}
		// Create event object
		e := Event{
			Name:  item.Name,
			Value: string(b),
			Order: order,
		}
		// Send event on channel
		c <- e
		// Wait set interval for next iteration
		// or continue right away if remote has been
		// triggered
		select {
		case <-ticker.C:
			continue
		case <-remote:
			fmt.Printf("[Worker %s] Remote was triggered\n", item.Name)
			continue
		}
	}
}

// RunCommand in bash using -c
func RunCommand(command string) ([]byte, error) {
	cmd := exec.Command("bash", "-c", command)
	return cmd.Output()
}

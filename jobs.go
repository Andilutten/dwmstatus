package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func JobTime() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04") + " "
}

func JobBattery() string {
	cmd := exec.Command("acpi")
	buf := new(bytes.Buffer)
	b, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[JobBattery] Error: %v\n", err)
		return ""
	}
	re := regexp.MustCompile(`(?m)Battery [0-9]: [a-zA-Z]*, ([0-9]{1,3}%)`)
	for _, match := range re.FindAllStringSubmatch(string(b), -1) {
		fmt.Fprintf(buf, " %s ", match[1])
	}
	return buf.String()
}

func JobWifi() string {
	cmd := exec.Command("iwgetid", "--raw")
	b, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[JobWifi] Error: %v\n", err)
		return ""
	}
	return strings.TrimSpace(string(b)) + " "
}

func JobVolume() string {
	args := strings.Split("awk -F\"[][]\" '/dB/ { print $2 }' <(amixer sget Master)"," ")
	cmd := exec.Command(args[0], args[1:]...)
	b, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[JobVolume] Error: %v\n", err)
		return ""
	}
	return "墳 " + strings.TrimSpace(string(b)) + " "
}

func JobBrightness() string {
	return ""
}

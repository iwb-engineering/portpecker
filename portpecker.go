package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/fatih/color"
)

type Port struct {
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

type Rule struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Ports       []Port `json:"ports"`
	Note        string `json:"note"`
	ApplyToAll  bool
}

type Config struct {
	Rules []Rule `json:"rules"`
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no suitable IP address found")
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Set ApplyToAll for rules with "*" as source
	for i := range config.Rules {
		if config.Rules[i].Source == "*" {
			config.Rules[i].ApplyToAll = true
		}
	}

	return &config, nil
}

func checkConnection(destination, protocol, port string) bool {
	address := net.JoinHostPort(destination, port)
	timeout := 3 * time.Second

	var conn net.Conn
	var err error

	if protocol == "TCP" {
		conn, err = net.DialTimeout("tcp", address, timeout)
	} else if protocol == "UDP" {
		conn, err = net.DialTimeout("udp", address, timeout)
	} else {
		return false
	}

	if err != nil {
		return false
	}

	conn.Close()
	return true
}

func main() {
	configFile := "config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	localIP, err := getLocalIP()
	if err != nil {
		color.Red("Error getting local IP: %v", err)
		os.Exit(1)
	}
	color.Cyan("Local IP address: %s", localIP)
	fmt.Println()

	config, err := loadConfig(configFile)
	if err != nil {
		color.Red("Error loading config: %v", err)
		os.Exit(1)
	}

	rulesChecked := false

	for _, rule := range config.Rules {
		if rule.Source == localIP || rule.ApplyToAll {
			rulesChecked = true
			color.Yellow("==== Checking rule: %s ====", rule.Note)
			if rule.ApplyToAll {
				color.Yellow("(This rule applies to all hosts)")
			}
			for _, port := range rule.Ports {
				success := checkConnection(rule.Destination, port.Protocol, port.Port)
				if success {
					color.Green("✓ %s:%s (%s) - Connection successful", rule.Destination, port.Port, port.Protocol)
				} else {
					color.Red("✗ %s:%s (%s) - Connection failed", rule.Destination, port.Port, port.Protocol)
				}
			}
			fmt.Println()
		}
	}

	if !rulesChecked {
		color.Red("No applicable rules found for this host")
	}
}

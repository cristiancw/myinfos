package info

import (
	"log"
	"net"
	"os"
	"time"
)

// Machine group informations about the machine and operational system.
type Machine struct {
	IPAddress string `json:"ip_address"`
	Hostname  string `json:"hostname"`
	Uptime    int64  `json:"uptime"`
	LastPing  int64  `json:"last_ping"`
}

// LoadMachine load and keep loading every 5 seconds to update the informations about it.
func LoadMachine(startTime time.Time) {
	for {
		machine := Machine{
			IPAddress: GetLocalIP(),
			Hostname:  getHostname(),
			Uptime:    time.Since(startTime).Nanoseconds() / 1000000000, // In seconds
			LastPing:  time.Now().Unix(),                                // In seconds
		}

		if err := SaveMachine(machine); err != nil {
			log.Fatal(err)
		}

		log.Printf("Refresh the machine info: %v\n", machine)
		time.Sleep(5 * time.Second)
	}
}

// GetLocalIP get local ip address.
func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80") // google address
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func getHostname() string {
	osName, err := os.Hostname()
	if err != nil {
		osName = "Unknown"
		log.Fatal(err)
	}
	return osName
}

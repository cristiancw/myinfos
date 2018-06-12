package info

import (
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

// Machine group informations about the machine and operational system.
type Machine struct {
	IPAddress    string `json:"ip_address"`
	Hostname     string `json:"hostname"`
	Uptime       int64  `json:"uptime"`
	RunningSince int64  `json:"running_since"`
}

// LoadMachine load and keep loading every 5 seconds to update the informations about it.
func LoadMachine(startTime time.Time) {
	for {
		machine := Machine{
			IPAddress:    getIP(),
			Hostname:     getHostname(),
			Uptime:       getUptime(),                                      // In seconds
			RunningSince: time.Since(startTime).Nanoseconds() / 1000000000, // In seconds
		}

		if err := SaveMachine(machine); err != nil {
			log.Fatal(err)
		}

		time.Sleep(5 * time.Second)
	}
}

func getIP() string {
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

func getUptime() int64 {
	info := syscall.Sysinfo_t{}
	syscall.Sysinfo(&info)
	return info.Uptime
}

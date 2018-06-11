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
	hostname     string
	ipAddress    string
	uptime       int64
	runningSince int64
}

// LoadMachine load and keep loading every 5 seconds to update the informations about it.
func LoadMachine(startTime time.Time, timeChan chan<- Machine) {
	for {
		machine := Machine{
			hostname:     getHostname(),
			ipAddress:    getIP(),
			uptime:       getUptime(),
			runningSince: time.Since(startTime).Nanoseconds(),
		}

		getIP()
		timeChan <- machine

		time.Sleep(5 * time.Second)
	}
}

func getHostname() string {
	osName, err := os.Hostname()
	if err != nil {
		osName = "Unknown"
		log.Fatal(err)
	}
	return osName
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

func getUptime() int64 {
	info := syscall.Sysinfo_t{}
	syscall.Sysinfo(&info)
	return info.Uptime
}

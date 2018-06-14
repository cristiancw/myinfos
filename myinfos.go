package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"./info"

	"github.com/gorilla/mux"
)

// Header to help to define the headers.
type Header map[string]string

func main() {
	go info.LoadMachine(time.Now())

	startServer()
}

func startServer() {
	router := mux.NewRouter()
	// Rest endpoints
	router.HandleFunc("/myinfos", getMachines).Methods("GET")

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8888", nil))
}

func getMachines(w http.ResponseWriter, r *http.Request) {
	machines, err := info.GetMachines()
	if err != nil {
		log.Fatal(err)
	}

	json, errJ := json.Marshal(machines)
	if errJ != nil {
		log.Fatal(err)
		response(w, http.StatusInternalServerError, nil, "")
	} else {
		accept := r.Header.Get("Accept")
		if strings.Contains(accept, "text/html") {
			html := createHTMLPage(machines)
			response(w, http.StatusOK, Header{"Content-Type": "text/html"}, html)
		} else {
			response(w, http.StatusOK, Header{"Content-Type": "application/json"}, string(json))
		}
	}
}

func response(w http.ResponseWriter, status int, header Header, body string) {
	for k, v := range header {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
	if body != "" {
		fmt.Fprintf(w, body)
	}
}

func createHTMLPage(machines []info.Machine) string {
	var buffer bytes.Buffer

	buffer.WriteString("<!DOCTYPE html>")
	buffer.WriteString("<html>")
	buffer.WriteString("<head>")
	buffer.WriteString("    <meta charset=\"UTF-8\"/>")
	buffer.WriteString("    <title>Document</title>")
	buffer.WriteString("</head>")
	buffer.WriteString("<body><center>")
	buffer.WriteString("    <h1>Machines</h1>")
	buffer.WriteString(table(machines))
	buffer.WriteString("</center></body>")
	buffer.WriteString("</html>")

	return buffer.String()
}

func table(machines []info.Machine) string {
	var buffer bytes.Buffer

	buffer.WriteString("<table cellspacing=\"10\"><tr>")
	i := 0
	for _, machine := range machines {
		online := time.Now().Unix() - machine.LastPing
		if online < 0 {
			online *= -1
		}

		i++
		if online > 10 {
			buffer.WriteString("<td style=\"padding:10px; margin:10px; background:#ffcccc;\">")
		} else {
			buffer.WriteString("<td style=\"padding:10px; margin:10px; background:#ccffdd;\">")
		}
		buffer.WriteString(machine.Hostname)
		buffer.WriteString("<hr>")
		buffer.WriteString("Ip address: ")
		buffer.WriteString(machine.IPAddress)
		buffer.WriteString("<br>OS uptime: ")
		buffer.WriteString(formatTime(machine.Uptime))
		buffer.WriteString("<br>Run since: ")
		buffer.WriteString(formatTime(machine.RunningSince))
		buffer.WriteString("<br>Last Ping: ")
		buffer.WriteString(fmt.Sprintf("%d s", online))
		buffer.WriteString("</td>")

		if i > 5 {
			i = 0
			buffer.WriteString("</tr><tr>")
		}
	}
	buffer.WriteString("</tr></table>")

	return buffer.String()
}

func formatTime(input int64) string {
	days := math.Floor(float64(input) / 60 / 60 / 24)
	seconds := input % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = input % 60

	var buffer bytes.Buffer

	if days > 0 {
		buffer.WriteString(fmt.Sprintf("%g day(s) ", days))
	}

	if hours < 10 {
		buffer.WriteString(fmt.Sprintf("0%g", hours))
	} else {
		buffer.WriteString(fmt.Sprintf("%g", hours))
	}

	if minutes < 10 {
		buffer.WriteString(fmt.Sprintf(":0%g", minutes))
	} else {
		buffer.WriteString(fmt.Sprintf(":%g", minutes))
	}

	if seconds < 10 {
		buffer.WriteString(fmt.Sprintf(":0%d", seconds))
	} else {
		buffer.WriteString(fmt.Sprintf(":%d", seconds))
	}

	return buffer.String()
}
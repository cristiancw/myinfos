package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cristiancw/myinfos/info"
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
	router.HandleFunc("/myinfos/machines", getMachines).Methods("GET")

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
		response(w, http.StatusOK, Header{"Content-Type": "application/json"}, string(json))
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

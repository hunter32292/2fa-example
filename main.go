package main

import (
	"log"
	"net/http"
)

const (
	// Username - Global Username for testing
	Username = "John"
	// Password - Global Password for testing
	Password = "IamNotACat"
)

func main() {
	log.Println("Staring...")
	http.HandleFunc("/login", login)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func login(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", r)
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(nil)
		return
	}

	if check(username, password) {
		w.WriteHeader(http.StatusAccepted)
		w.Write(nil)
		return
	}

	w.WriteHeader(http.StatusForbidden)
	w.Write(nil)
	return
}

// check - Check the username and password against some sort of authentication data store
func check(username, password string) bool {
	if username == Username && password == Password {
		return true
	}
	return false
}

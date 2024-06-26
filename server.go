package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	Username string `json:"user"`
	Password string `json:"password"`
}

var user User
var PORT = ":1234"
var DATA = make(map[string]string)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusNotFound)
	Body := "Thanks for visiting!"
	fmt.Fprintf(w, "%s\n", Body)
}



func timeHandler (w http.ResponseWriter, r *http.Request) {
	log.Println("Serving", r.URL.Path, "from", r.Host)
	t := time.Now().Format(time.RFC1123)
	Body := "The current time is" + t + "\n"
	fmt.Fprintf(w, "%s", Body)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving", r.URL.Path, "from", r.Host, r.Method)

	if r.Method != http.MethodPost {
		http.Error(w, "Error", http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%s\n", "Method not allowed")
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(d, &user)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}

	if user.Username != "" {
		DATA[user.Username] = user.Password
		log.Println(DATA)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving", r.URL.Path, "from", r.Host, r.Method)

	if r.Method != http.MethodGet {
		http.Error(w, "Error", http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%s\n", "Method not allowed")
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Readall Error", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(d, &user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unmarshal error", http.StatusBadRequest)
		return
	}

	fmt.Println(user)

	_, ok := DATA[user.Username]
	if ok && user.Username != "" {
		log.Println("Found!")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s\n", d)
	} else {
		log.Println("Not found!")
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, "Map - resource not found", http.StatusNotFound)
		return
	}

}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host, r.Method)
	if r.Method != http.MethodDelete {
		http.Error(w, "error", http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%s\n", "Method not allowd")
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Readall - error", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(d, &user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unmarshal error", http.StatusBadRequest)
		return
	}

	log.Println(user)

	_, ok := DATA[user.Username]
	if ok && user.Username != "" {
		if user.Password == DATA[user.Username] {
			delete(DATA, user.Username)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s\n", d)
			log.Println(DATA)
		}
	} else {
		log.Println("User", user.Username, "not found!")
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, "Delete - resource not found", http.StatusNotFound)
		return
	}
	log.Println("After:", DATA)
}

func main() {

	arguments := os.Args

	if len(arguments) != 1 {
		PORT = ":" + arguments[1]
	}

	mux := http.NewServeMux()
	s := &http.Server{
		Addr: PORT,
		Handler: mux,
		IdleTimeout: 10 * time.Second,
		ReadTimeout: time.Second,
		WriteTimeout: time.Second,
	}


	mux.Handle("/time", http.HandlerFunc(timeHandler))
	mux.Handle("/add", http.HandlerFunc(addHandler))
	mux.Handle("/delete", http.HandlerFunc(deleteHandler))
	mux.Handle("/", http.HandlerFunc(defaultHandler))
	mux.Handle("/get", http.HandlerFunc(getHandler))


	fmt.Println("Ready to serve at", PORT)
	err := s.ListenAndServe()

	if err != nil{
		fmt.Println(err)
		return
	}
}

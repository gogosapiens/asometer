package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

func main() {
	http.HandleFunc("/measure", measureHandler)
	http.HandleFunc("/progress", progressHandler)
    http.ListenAndServe(":8080", nil)
}

func progressHandler(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    fmt.Fprintf(w, fmt.Sprintf("%s", getProgress(q.Get("title"))))
    w.Header().Set("Content-Type", "application/json")
   
   	json, _ := json.Marshal(progress.GetStates())
    w.Write([]byte(json))
    //w.WriteHeader(http.StatusOK)
}

func measureHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
    q := r.URL.Query()
    go measureTraffic(q.Get("title"), q.Get("keywords"), q.Get("country"))
    w.WriteHeader(http.StatusOK)
}
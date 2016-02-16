package main

import (
	"net/http"
	//"fmt"
	"encoding/json"
)

func main() {
    http.HandleFunc("/measure", measureHandler)
    http.HandleFunc("/progress", progressHandler)
    http.ListenAndServe(":8080", nil)
}

func progressHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
   
   	json, _ := json.Marshal(progress.GetStates())
    w.Write([]byte(json))
}

func measureHandler(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    go measureTraffic(q.Get("title"), q.Get("keywords"), q.Get("country"), q.Get("email"))
    w.WriteHeader(http.StatusOK)
}
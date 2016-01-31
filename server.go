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
  //  	json := `[{"Name":"Adam","Age":36,"Job":"CEO"},
  // {"Name":"Eve","Age":34,"Job":"CFO"},
  // {"Name":"Mike","Age":38,"Job":"COO"}]`
    w.Write([]byte(json))
    //w.WriteHeader(http.StatusOK)
}

func measureHandler(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    go measureTraffic(q.Get("title"), q.Get("keywords"), q.Get("country"))
    w.WriteHeader(http.StatusOK)
}
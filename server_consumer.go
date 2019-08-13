package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
)

//StockSchema Used for stock schema
type StockSchema struct {
	Symbol string
	Price  float64
}

//StockBatch Used for stock schema
type StockBatch struct {
	Items []StockSchema
}

var totalRequestProcessed = 0

func main() {
	var port = ":3001"
	setUpRoutes()
	fmt.Println("Http server running on port", port)
	http.ListenAndServe(port, nil)
}

func setUpRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/stock/", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r)
	})
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg StockBatch
	err = json.Unmarshal(b, &msg)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//log.Println(string(output))
	w.Header().Set("content-type", "application/json")
	w.Write(output)
	totalRequestProcessed++
	if totalRequestProcessed%4000 == 0 {
		log.Printf("Request Processed: %d\n", totalRequestProcessed)
	}
}

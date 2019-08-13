package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var stockList = []string{"msft", "goog", "amzn", "ntnx", "netflix", "fb", "uber", "apl"}
var priceList = []float64{140.0, 1200, 1950, 20, 340, 190, 40, 200}

func main() {
	start := time.Now()
	var totalReq = 10000
	for w := 1; w <= totalReq; w++ {
		sendPostRequest(getNextStockUpdate())
	}

	elapsed := time.Since(start)
	log.Printf("Time taken %s for %d requests", elapsed, totalReq)
	//test()
}

func sendStocksPostRequest(stock string, price float64) {
	url := "http://localhost:3001/stock/"
	message := map[string]interface{}{
		"Symbol": stock,
		"Price":  price,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	//log.Println(result)
}

func sendPostRequest(data []byte) {
	url := "http://localhost:3001/stock/"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	//log.Println(result)
}

func sendPostRequestRespOnChannel(data []byte, results chan<- []byte) {
	sendPostRequest(data)
	results <- data
}

func worker(id int, jobs <-chan []byte, results chan<- []byte) {
	for j := range jobs {
		go sendPostRequestRespOnChannel(j, results)
	}
}

func test() {
	start := time.Now()
	var nWorkers = 10
	// In order to use our pool of workers we need to send
	// them work and collect their results. We make 2
	// channels for this.
	jobs := make(chan []byte, nWorkers)
	results := make(chan []byte, nWorkers)

	// This starts up 3 workers, initially blocked
	// because there are no jobs yet.

	for w := 1; w <= nWorkers; w++ {
		go worker(w, jobs, results)
	}

	var totalReq = 10000
	var interations = totalReq / nWorkers
	for j := 1; j <= interations; j++ {
		for a := 1; a <= nWorkers; a++ {
			jobs <- getNextStockUpdate()
		}

		for a := 1; a <= nWorkers; a++ {
			<-results
		}
	}

	close(jobs)
	elapsed := time.Since(start)
	log.Printf("Time taken %s for %d requests", elapsed, totalReq)
}

func getNextStockUpdate() []byte {
	var idx = rand.Intn(8)
	var stock = stockList[idx]
	var price = priceList[idx] + rand.Float64()
	message := map[string]interface{}{
		"Symbol": stock,
		"Price":  price,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return nil
	}

	return bytesRepresentation
}

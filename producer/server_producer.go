package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
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

var stockList = []string{"msft", "goog", "amzn", "ntnx", "netflix", "fb", "uber", "apl"}
var priceList = []float64{140.0, 1200, 1950, 20, 340, 190, 40, 200}

var url = "http://139.59.88.85:3001/stock/"

//var url = "http://localhost:3001/stock/"

func main() {
	start := time.Now()
	var totalRequest = 50000
	testWithWorker(totalRequest)
	//testSimple(totalRequest)
	elapsed := time.Since(start)
	log.Printf("Time taken %s for %d requests", elapsed, totalRequest)
}

func sendPostRequest(data []byte) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatalln(err)
	}

	decodeResponse(resp)
}

func decodeResponse(resp *http.Response) {
	//b, err := ioutil.ReadAll(resp.Body)
	ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	//if err != nil {
	//	return
	//}

	//log.Println(string(b))
	//var msg StockBatch
	//err = json.Unmarshal(b, &msg)
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
func testSimple(totalReq int) {
	for w := 1; w <= totalReq; w++ {
		sendPostRequest(getNextStockUpdate())
	}
}

func testWithWorker(totalReq int) {
	var nWorkers = 50
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

	var interations = totalReq / nWorkers
	for j := 1; j <= interations; j++ {
		for a := 1; a <= nWorkers; a++ {
			//jobs <- getNextStockUpdate()
			jobs <- getNextStockUpdateBatch(100)
		}

		for a := 1; a <= nWorkers; a++ {
			<-results
		}
	}

	close(jobs)
}

func getNextStockUpdate() []byte {
	message := getNextStockUpdateInSchema()
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return nil
	}

	return bytesRepresentation
}

func getNextStockUpdateBatch(batchSize int) []byte {
	var batch StockBatch
	for kk := 1; kk <= batchSize; kk++ {
		var item = getNextStockUpdateInSchema()
		batch.Items = append(batch.Items, item)
	}

	bytes, err := json.Marshal(batch)
	if err != nil {
		log.Println(err)
		return nil
	}

	return bytes
}

func getNextStockUpdateInSchema() StockSchema {
	var idx = rand.Intn(8)
	var item StockSchema
	item.Symbol = stockList[idx]
	item.Price = priceList[idx] + rand.Float64()
	return item
}

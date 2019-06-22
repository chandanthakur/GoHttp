package main

import (
	"fmt"
	"net/http"
)

func main() {
	var port = ":3000"
	http.Handle("/", http.FileServer(http.Dir("./gerritstats")))
	fmt.Println("Http server running on port", port)
	http.ListenAndServe(port, nil)
}

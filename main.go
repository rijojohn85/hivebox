package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rijojohn85/hivebox/hivebox"
)

func main() {
	http.HandleFunc("/version", hivebox.GetVersion)
	http.HandleFunc("/temperature", hivebox.GetAvgTemp)

	fmt.Println("Server is running at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/arpradhan/usda-calorie-search"
)

func main() {
	queryPtr := flag.String("q", "", "A food query")
	apiKeyPtr := flag.String("apikey", "", "API Key for ndb.nal.usda.gov")

	flag.Parse()

	if *queryPtr == "" {
		log.Fatal("query is required.")
	}

	if *apiKeyPtr == "" {
		log.Fatal("apikey is required.")
	}

	client := usda.NewCalorieSearchClient(*apiKeyPtr)
	calorieResponse, err := client.Get(*queryPtr)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(calorieResponse)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

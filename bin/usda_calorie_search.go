package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/arpradhan/usda-calorie-search"
)

func main() {
	queryPtr := flag.String("query", "", "A food query")
	apiKeyPtr := flag.String("apikey", "", "API Key for ndb.nal.usda.gov")

	flag.Parse()

	if *queryPtr == "" {
		log.Fatal("query is required.")
	}

	if *apiKeyPtr == "" {
		log.Fatal("apikey is required.")
	}

	searchClient := usda.NewUSDASearchClient(*apiKeyPtr)

	resp, err := searchClient.Get(*queryPtr)

	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal(fmt.Errorf("\nURL: %s \nStatus: %s \nBody: %s", resp.Request.URL, resp.Status, b))
	}

	searchResponse := new(usda.SearchResponse)

	err = json.Unmarshal(b, &searchResponse)
	if err != nil {
		log.Fatal(err)
	}

	if len(searchResponse.List.Item) == 0 {
		log.Fatal(fmt.Errorf("Food not found"))
	}

	calorieResponse := new(usda.CalorieResponse)

	nutrientClient := usda.NewUSDANutrientClient(*apiKeyPtr)

	for _, item := range searchResponse.List.Item {
		resp, err = nutrientClient.Get(item.Ndbno)

		if err != nil {
			log.Fatal(err)
		}

		b, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode != 200 {
			log.Fatal(fmt.Errorf("\nURL: %s \nStatus: %s \nBody: %s", resp.Request.URL, resp.Status, b))
		}

		nutrientResponse := new(usda.NutrientResponse)

		err = json.Unmarshal(b, &nutrientResponse)
		if err != nil {
			log.Fatal(err)
		}

		if len(nutrientResponse.Report.Foods) > 0 {
			calorieResponse.Foods = append(calorieResponse.Foods, nutrientResponse.Report.Foods[0])
		}
	}

	b, err = json.Marshal(calorieResponse)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

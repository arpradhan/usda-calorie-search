package usda

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type SearchResponse struct {
	List SearchList `json:"list"`
}

type SearchList struct {
	Q     string       `json:"q"`
	Sr    string       `json:"sr"`
	Ds    string       `json:"ds"`
	Start int          `json:"start"`
	End   int          `json:"end"`
	Total int          `json:"total"`
	Group string       `json:"group"`
	Sort  string       `json:"sort"`
	Item  []SearchItem `json:"item"`
}

type SearchItem struct {
	Offset int    `json:"offset"`
	Group  string `json:"group"`
	Name   string `json:"name"`
	Ndbno  string `json:"ndbno"`
	Ds     string `json:"ds"`
	Manu   string `json:"manu"`
}

type NutrientResponse struct {
	Report Report `json:"report"`
}

type Report struct {
	Sr     string `json:"sr"`
	Groups string `json:"groups"`
	Subset string `json:"subset"`
	End    int    `json:"end"`
	Start  int    `json:"start"`
	Total  int    `json:"total"`
	Foods  []Food `json:"foods"`
}

type Food struct {
	Ndbno     string     `json:"ndbno"`
	Name      string     `json:"name"`
	Weight    float32    `json:"weight"`
	Measure   string     `json:"measure"`
	Nutrients []Nutrient `json:"nutrients"`
}

type Nutrient struct {
	NutrientID string  `json:"nutrient_id"`
	Nutrient   string  `json:"nutrient"`
	Unit       string  `json:"unit"`
	Value      string  `json:"value"`
	Gm         float32 `json:"gm"`
}

type CalorieResponse struct {
	Foods []Food `json:"foods"`
}

type USDASearchClient struct {
	URL    string
	Format string
	Sort   string
	Max    int
	Offset int
	APIKey string
}

type USDANutrientClient struct {
	URL       string
	Format    string
	Nutrients string
	Max       int
	APIKey    string
}

type CalorieSearchClient struct {
	USDASearchClient   *USDASearchClient
	USDANutrientClient *USDANutrientClient
}

func NewCalorieSearchClient(apiKey string) *CalorieSearchClient {
	return &CalorieSearchClient{
		USDASearchClient:   NewUSDASearchClient(apiKey),
		USDANutrientClient: NewUSDANutrientClient(apiKey),
	}
}

func NewUSDASearchClient(apiKey string) *USDASearchClient {
	return &USDASearchClient{
		URL:    "https://api.nal.usda.gov/ndb/search/",
		Format: "json",
		Sort:   "n",
		Max:    5,
		Offset: 0,
		APIKey: apiKey,
	}
}

func NewUSDANutrientClient(apiKey string) *USDANutrientClient {
	return &USDANutrientClient{
		URL:       "https://api.nal.usda.gov/ndb/nutrients/",
		Format:    "json",
		Nutrients: "208",
		Max:       1,
		APIKey:    apiKey,
	}
}

func (client *USDANutrientClient) Get(ndbNo string) (*http.Response, error) {
	url := fmt.Sprintf(
		"%v?format=%v&nutrients=%v&max=%v&api_key=%v&ndbno=%v",
		client.URL,
		client.Format,
		client.Nutrients,
		client.Max,
		client.APIKey,
		ndbNo,
	)
	return http.Get(url)
}

func (client *USDASearchClient) Get(query string) (*http.Response, error) {
	query = strings.Replace(query, " ", "%20", -1)
	url := fmt.Sprintf(
		"%v?format=%v&sort=%v&max=%v&offset=%v&api_key=%v&q=%v",
		client.URL,
		client.Format,
		client.Sort,
		client.Max,
		client.Offset,
		client.APIKey,
		query,
	)
	return http.Get(url)
}

func (client *CalorieSearchClient) Get(query string) (*CalorieResponse, error) {
	resp, err := client.USDASearchClient.Get(query)
	calorieResponse := new(CalorieResponse)

	if err != nil {
		return calorieResponse, err
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return calorieResponse, err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("\nURL: %s \nStatus: %s \nBody: %s", resp.Request.URL, resp.Status, b)
		return calorieResponse, err
	}

	searchResponse := new(SearchResponse)

	err = json.Unmarshal(b, &searchResponse)
	if err != nil {
		return calorieResponse, err
	}

	if len(searchResponse.List.Item) == 0 {
		log.Fatal(fmt.Errorf("Food not found"))
	}

	for _, item := range searchResponse.List.Item {
		resp, err = client.USDANutrientClient.Get(item.Ndbno)

		if err != nil {
			return calorieResponse, err
		}

		b, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			return calorieResponse, err
		}

		if resp.StatusCode != 200 {
			err = fmt.Errorf("\nURL: %s \nStatus: %s \nBody: %s", resp.Request.URL, resp.Status, b)
			return calorieResponse, err
		}

		nutrientResponse := new(NutrientResponse)

		err = json.Unmarshal(b, &nutrientResponse)
		if err != nil {
			return calorieResponse, err
		}

		if len(nutrientResponse.Report.Foods) > 0 {
			calorieResponse.Foods = append(calorieResponse.Foods, nutrientResponse.Report.Foods[0])
		}
	}
	return calorieResponse, nil
}

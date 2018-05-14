package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/arpradhan/usda-calorie-search"
	"github.com/kelseyhightower/envconfig"
)

type Specification struct {
	APIKey string
}

func NewSpecification() *Specification {
	var s Specification
	err := envconfig.Process("usda", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &s
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	spec := NewSpecification()
	r.ParseForm()
	q := r.FormValue("q")
	if q == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "q parameter is required.")
	} else {
		client := usda.NewCalorieSearchClient(spec.APIKey)
		calorieResponse, err := client.Get(q)

		if err != nil {
			log.Fatal(err)
		}

		b, err := json.Marshal(calorieResponse)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(w, string(b))
	}

}

func main() {
	http.HandleFunc("/api/v1/foods", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

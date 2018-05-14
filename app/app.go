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
		log.Printf("%v %v %v", r.Method, r.URL, http.StatusBadRequest)
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
		log.Printf("%v %v %v", r.Method, r.URL, http.StatusOK)
		fmt.Fprintf(w, string(b))
	}

}

func main() {
	http.HandleFunc("/api/v1/foods", handler)
	addr := ":8080"
	log.Printf("Running app on %v", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

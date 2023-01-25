package main

import (
	"fmt"
	"html/template"
	"lab6/knna"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var rek *knna.Knn
var learningData []knna.Data
var pageData []knna.Data

func itemFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)               // get variables from router as map
	id, _ := strconv.Atoi(vars["id"]) // get id value from map
	fmt.Println("Id: ", id)           // print id value
	tmpl, _ := template.ParseFiles("pages/apt.html")
	apt := pageData[id]
	apt.Price = rek.Predict(apt.X, 5)
	tmpl.Execute(w, apt)
}

//adres http://localhost:8080/apts/medium

func aptsFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sizeType := vars["jakie"]
	var smallApts, mediumApts, largeApts []knna.Data

	for _, apt := range pageData {
		if sizeType == "small" && apt.Area < 50 {
			smallApts = append(smallApts, apt)
		} else if sizeType == "large" && apt.Area > 100 {
			largeApts = append(largeApts, apt)
		} else if sizeType == "medium" && apt.Area >= 50 && apt.Area <= 100 {
			mediumApts = append(mediumApts, apt)
		}
	}

	tmpl, _ := template.ParseFiles("pages/apts.html")
	tmpl.Execute(w, map[string]interface{}{
		"SmallApts":  smallApts,
		"MediumApts": mediumApts,
		"LargeApts":  largeApts,
	})
}

func main() {
	// wczytanie danych i utworzenie systemu
	cols := 10
	learningData = knna.LoadData("data/apt-use.txt", cols)
	pageData = knna.LoadData("data/apt-test.txt", cols)
	rek = knna.NewKnn(learningData)
	// test systemu dla próbki od id równym zero, oraz k = 5
	results_from_system := rek.Predict(pageData[0].X, 5)
	results_from_data := pageData[0].Y
	fmt.Printf("APARTAMENT[0]: %+v\n", pageData[0])
	fmt.Println("SYSTEM: ", results_from_system)
	fmt.Println("OGLOSZENIE: ", results_from_data)
	// zadanie 1 - obliczenie MAE dla wszystkich próbek pageData
	start := time.Now()
	var sum float64
	for i := range pageData {
		predicted := rek.Predict(pageData[i].X, 5)
		diff := math.Abs(predicted - pageData[i].Y)
		sum += diff
	}
	mae := sum / float64(len(pageData))
	mae = mae / float64(len(pageData))
	fmt.Println("MAE: ", mae)
	elapsed := time.Since(start)
	fmt.Println("Czas potrzebny na obliczenie MAE: ", elapsed)

	router := mux.NewRouter()
	router.HandleFunc("/apt/{id}", itemFunc).Methods("GET")
	router.HandleFunc("/apts/{jakie}", aptsFunc).Methods("GET")
	http.ListenAndServe(":8080", router)

}

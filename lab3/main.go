package main

import (
	"fmt"
	"html/template"
	"lab0302/knn"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var recsys *knn.Knn
func piwaFunc(w http.ResponseWriter, r *http.Request) {
	// pobranie 10 losowych piw z systemu
	tenbeers := recsys.Get10RandomBeers()
	// przekazanie 10 losowych piw do strony pages/beers.html
	tmpl, _ := template.ParseFiles("pages/beer.html")
	tmpl.Execute(w, tenbeers)
   } 
   func rekoFunc(w http.ResponseWriter, r *http.Request) {
	// gdy metoda inna niż POST - error 403
	if r.Method != "POST" {
	http.Error(w, "Only POST supported", http.StatusForbidden)
	return
	}
	// parsowanie danych z formularza
	r.ParseForm()
	// iteracja danych z formularza i jednoczesna ocena piw
	for name, element := range r.PostForm {
	rate, _ := strconv.ParseFloat(element[0], 64) // parsowanie
	id, _ := strconv.Atoi(name[4:]) // wycięcie id
	beer := recsys.GetBeerByID(id) // pobranie piwa
	beer.Rate = rate // ocena piwa
	fmt.Println("Name:", name, ", Id:", id, ", Rate:", rate)
	}
	// pobranie rekomendacji
	recbeers := recsys.GetRecommendation()
	tmpl, _ := template.ParseFiles("pages/reko.html")
	tmpl.Execute(w, recbeers)
   } 
   
func main() {

	// po testach działania nalezy te linijki zostawic
	rand.Seed(time.Now().UnixMilli())
	recsys = knn.Initialize()

	http.HandleFunc("/piwa/", piwaFunc)
	http.HandleFunc("/reko/", rekoFunc)
	http.ListenAndServe("localhost:8080", nil)
	// a te ponizej zakomentowac
	 rand.Seed(time.Now().UnixMilli())
	 recsys = knn.Initialize()
	 var rateText string
	 var rateCount int = 0
	 for rateCount < 10 {
		beer := recsys.GetRandomBeer()
	 	beer.DisplayInformation(recsys)
	 	fmt.Println("Ocen piwo (1-5): ")
	 	fmt.Scanln(&rateText)
	 	rate, error := strconv.ParseFloat(rateText, 64)
	 	if error == nil && rate >= 1 && rate <= 5 {
	 		beer.Rate = rate
	 		rateCount++
	 	}
	 }
	 reco := recsys.GetRecommendation()
	 fmt.Println("Rekomendowane piwa:")
	 reco[0].DisplayInformation(recsys)
	 reco[1].DisplayInformation(recsys)
	 reco[2].DisplayInformation(recsys)
}

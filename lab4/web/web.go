package web

import (
	"encoding/json"
	"fmt"
	"lab0302/knn"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var recsys *knn.Knn

type Beersandstyle struct {
	Beers []*knn.Beer
	Style *knn.Styles
}

func getID(uri string) int {
	data := strings.Split(uri, "/")
	last := data[len(data)-1]
	id, err := strconv.Atoi(last)
	if err == nil {
		return id
	}
	return -1
}

func indexFunc(w http.ResponseWriter, r *http.Request) {
	tenbeers := recsys.Get12RandomBeers()
	senddata := Beersandstyle{tenbeers, recsys.GetStyles()}
	tmpl, _ := template.ParseFiles("pages/index.html")
	tmpl.Execute(w, senddata)
}

func recoFunc(w http.ResponseWriter, r *http.Request) {
	data := strings.Split(r.RequestURI, "/")
	if len(data) < 3 {
		fmt.Fprintf(w, "ERROR - wrong URL")
		return
	}
	pairs := strings.Split(data[2], ";")
	if len(pairs) < 6 {
		fmt.Fprintf(w, "ERROR - too less rates")
		return
	}
	for i := 0; i < len(pairs)-1; i += 2 {
		bid, _ := strconv.Atoi(pairs[i])
		brate, _ := strconv.Atoi(pairs[i+1])
		// fmt.Println(bid, brate)
		beer := recsys.GetBeerByID(bid)
		beer.Rate = float64(brate)
	}
	reco := recsys.GetRecommendation()
	senddata := Beersandstyle{reco, recsys.GetStyles()}
	tmpl, _ := template.ParseFiles("pages/reco.html")
	tmpl.Execute(w, senddata)
}

func explFunc(w http.ResponseWriter, r *http.Request) {
	// obtain beer ID
	beerid := getID(r.RequestURI)
	// obtain beer
	beer := recsys.GetBeerByID(beerid)
	fmt.Println(beer)
	// obtain 3 similiar beers
	if beer != nil {
		beers := recsys.GetSimiliar(beer)
		fmt.Println(beers)

		// zadanie 2 - część 1 - zwrocic beers jako JSON
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(beers)
		w.Write(data)
	} else {
		// zadanie 2 - część 2 - zwrocic informacje o zlym ID jako JSON
		beer := knn.Beer{
			Id: -1,
		Abv: 0.01,
		Ibu: 0.001,
		Style: '-',
		Name: "-",
		Rate: 0.001,
		Estim: 0.001,
		Simi:  []*knn.Beer{},
		}
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(beer)
		w.Write(data)
	}
}

func StartServer() {
	rand.Seed(time.Now().UnixMilli())
	recsys = knn.Initialize()

	fs := http.FileServer(http.Dir("./data"))
	http.Handle("/files/", http.StripPrefix("/files/", fs))
	http.HandleFunc("/", indexFunc)
	http.HandleFunc("/reco/", recoFunc)
	http.HandleFunc("/expl/", explFunc)
	http.ListenAndServe("localhost:8080", nil)
}

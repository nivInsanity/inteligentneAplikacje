package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"
)

type SO struct {
	Movie     int
	UserArate float64
	UserBrate float64
}

func (u1 *User) Compare(u2 *User) float64 {
	if (len(u1.Rates) == 0) || (len(u2.Rates) == 0) {
		return math.MaxFloat64
	}
	if u1.Id == u2.Id {
		return math.MaxFloat64
	}
	same := []SO{}
	for i := 0; i < len(u1.Rates); i++ {
		if u1.Rates[i].Rate == 0.0 {
			continue
		}
		for j := 0; j < len(u2.Rates); j++ {
			if u2.Rates[j].Rate == 0.0 {
				continue
			}
			if u1.Rates[i].Movie == u2.Rates[j].Movie {
				takisam := SO{}
				takisam.Movie = u1.Rates[i].Movie
				takisam.UserArate = u1.Rates[i].Rate
				takisam.UserBrate = u2.Rates[j].Rate
				same = append(same, takisam)
			}
		}
	}
	C := float64(len(same))
	D := float64(len(u1.Rates) + len(u2.Rates))
	A := C / D
	B := 0.0
	for i := 0; i < len(same); i++ {
		B += math.Abs(same[i].UserArate - same[i].UserBrate)
	}
	B /= float64(len(same))
	roznica := (0.5 + A) * B

	return roznica
}

func GetUser(host string, uid int) *User {

	userData := User{}
	client := http.Client{Timeout: time.Second}
	url := host + "Users/" + strconv.Itoa(uid)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	execute, _ := client.Do(request)
	body, _ := ioutil.ReadAll(execute.Body)
	json.Unmarshal(body, &userData)

	return &userData
}

func GetMovies(host string) []Movie {
	movieData := []Movie{}
	client := http.Client{Timeout: time.Second}
	url := host
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	execute, _ := client.Do(request)
	body, _ := ioutil.ReadAll(execute.Body)
	json.Unmarshal(body, &movieData)

	return movieData
}

type Rate struct {
	Movie int     `json:"movieid"`
	Rate  float64 `json:"rate"`
}

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Rates []Rate `json:"rates"`
}

type Movie struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Oryginal string `json:"oryginal"`
	Year     int    `json:"year"`
	Genre    string `json:"genre"`
}

func main() {
	var tak float64 = 100
	id := 1
	userData1 := GetUser("http://localhost:3000/", 3)
	filmy := GetMovies("http://localhost:3000/Movies")
	otherUser1 := GetUser("http://localhost:3000/", 1)
	fmt.Println("Moje dane: ", userData1)
	fmt.Println("Dane innego użytkownika: ", otherUser1)
	roznica := userData1.Compare(otherUser1)
	fmt.Print("Różnica: ", roznica, "\n ")
	fmt.Print("\n")

	for i := 0; i < 5; i++ {
		otherUser1 := GetUser("http://localhost:3000/", i)
		fmt.Print("Różnica JA kontra użytkownik ", i, " ")
		roznica := userData1.Compare(otherUser1)
		fmt.Print(roznica, "\n")
		if tak > roznica {
			tak = roznica
			id = i
		}
	}
	fmt.Print("\n")
	fmt.Print("Najbardziej podobny uzytkownik do mnie: ", id)
	fmt.Print("\n ")
	podobny := GetUser("http://localhost:3000/", id)
	fmt.Print("\n ")
	fmt.Print("Obiekty z wysokimi ocenami użytkownika podobnego do mnie:\n ")

	for i := 0; i < len(podobny.Rates); i++ {
		if podobny.Rates[i].Rate > 9 {
			fmt.Printf("%+v\n", filmy[podobny.Rates[i].Movie])
		}
	}
}

package knna

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Data struct {
	ID            int       // id
	X             []float64 // wektor parametrów x (do prognozy)
	Y             float64   // cena (do porównania)
	Year          float64   // rok budowy
	Age           float64   // wiek budynku
	Area          float64   // metraż (metry kwadratowe)
	Floor         int       // numer piętra
	Pcount        int       // liczba miejsc parkingowych
	BusDistance   float64   // odległosć do przystanku autobusoego
	MetroDistance float64   // odległość do metra
	Facilities    int       // liczba ważniejszych placówek w okolicy
	Parks         int       // liczba parków w okolicy
	Schools       int       // liczba szkół w okolicy
	Price         float64   // cena z ogłoszenia (w zł)
}

func NewData(id int, x []float64, y float64) *Data {
	data := Data{ID: id, X: x, Y: y}
	data.Year = x[0]
	data.Age = x[1]
	data.Area = x[2]
	data.Floor = int(x[3])
	data.Pcount = int(x[4])
	data.BusDistance = x[5]
	data.MetroDistance = x[6]
	data.Facilities = int(x[7])
	data.Parks = int(x[8])
	data.Schools = int(x[9])
	data.Price = y
	return &data
}

func LoadData(filename string, inputs int) []Data {
	id := 0
	data := []Data{}
	file, _ := os.Open(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "\t")
		x := make([]float64, inputs)
		for i := 0; i < inputs; i++ {
			x[i], _ = strconv.ParseFloat(line[i], 64)
		}
		y, _ := strconv.ParseFloat(line[inputs], 64)
		data = append(data, *NewData(id, x, y))
		id++
	}
	return data
}

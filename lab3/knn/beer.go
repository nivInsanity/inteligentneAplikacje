package knn

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Beer struct {
	Id    int     // id piwa
	Abv   float64 // poziom abv
	Ibu   float64 // goryczka
	Style int     // styl
	Name  string  // nazwa
	Rate  float64 // ocena
	Estim float64 // estymowana ocena
}

type Beers struct {
	beers  []Beer
	styles Styles
}

func (b *Beer) DisplayInformation(knn *Knn) {
	fmt.Print("Nazwa: ", b.Name)
	fmt.Print(", Alkohol: ", b.Abv*100.0, "%")
	fmt.Print(", Goryczka: ", b.Ibu)
	fmt.Print(", Styl: ", knn.GetStyles().GetStyleName(b.Style))
	if b.Estim > 0 {
		fmt.Print(", ERate: ", b.Estim)
	}
	fmt.Println()
}

func (b1 *Beer) Distance(b2 *Beer) float64 {
	var d float64 = 0
	if b1.Style != b2.Style {
		d += 1.0
	}
	d += math.Abs(b1.Abv-b2.Abv) * 5.0
	d += math.Abs(b1.Ibu-b2.Ibu) * 0.01
	return d
}

func LoadBeers(name string) *Beers {
	b := Beers{}
	file, error := os.Open(name)
	if error != nil {
		fmt.Println(error.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stringData := scanner.Text()
		data := strings.Split(stringData, ",")
		id, _ := strconv.Atoi(data[0])
		abv, _ := strconv.ParseFloat(data[1], 64)
		ibu, _ := strconv.ParseFloat(data[2], 64)
		style := data[5]
		name := data[4]
		if len(style) < 2 {
			continue
		}
		if ibu == 0.0 {
			continue
		}
		styleId := b.styles.CheckStyle(style)
		b.beers = append(b.beers, Beer{id, abv, ibu, styleId, name, 0, 0})
	}
	return &b
}

func (b *Beer) EstimateRate(bo []*Beer) {
	d1, d2, d3 := math.MaxFloat64, math.MaxFloat64, math.MaxFloat64
	i1, i2, i3 := 0, 0, 0
	for i := 0; i < len(bo); i++ {
		tmp := b.Distance(bo[i])
		if tmp < d1 {
			d3 = d2
			d2 = d1
			d1 = tmp
			i3 = i2
			i2 = i1
			i1 = i
		} else if tmp < d2 {
			d3 = d2
			d2 = tmp
			i3 = i2
			i2 = i
		} else if tmp < d3 {
			d3 = tmp
			i3 = i
		}
	}
	b.Estim = (bo[i1].Rate + bo[i2].Rate + bo[i3].Rate) / 3.0
}

func (b *Beers) Recomendation() []*Beer {
	// find rated beers
	beers := []*Beer{}
	for k, l := 0, len(b.beers); k < l; k++ {
		if b.beers[k].Rate > 0 {
			beers = append(beers, &b.beers[k])
		}
	}
	// calculate estimated rates
	for k, l := 0, len(b.beers); k < l; k++ {
		if b.beers[k].Rate == 0 {
			b.beers[k].EstimateRate(beers)
		}
	}
	// find best 3 results
	o1, o2, o3 := 0.0, 0.0, 0.0
	i1, i2, i3 := 0, 0, 0
	// o1 = b.beers[0].Estim
	for k, l := 0, len(b.beers); k < l; k++ {
		tmp := b.beers[k].Estim
		if tmp > o1 {
			o3 = o2
			o2 = o1
			o1 = tmp
			i3 = i2
			i2 = i1
			i1 = k
		} else if tmp > o2 {
			o3 = o2
			o2 = tmp
			i3 = i2
			i2 = k
		} else if tmp > o3 {
			o3 = tmp
			i3 = k
		}
	}
	// erase rates
	for k, l := 0, len(b.beers); k < l; k++ {
		if b.beers[k].Rate > 0 {
			b.beers[k].Rate = 0
		}
	}
	// return results
	var recom []*Beer = []*Beer{&b.beers[i1], &b.beers[i2], &b.beers[i3]}
	return recom
}

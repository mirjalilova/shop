package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Product struct {
	TovarNomi   string  `json:"tovar_nomi"`
	Razmer      int     `json:"razmer"`
	Soni        int     `json:"soni"`
	KelishNarxi float64 `json:"kelish_narxi"`
	SotishNarxi float64 `json:"sotish_narxi"`
}

func main() {
	file, err := os.Open("moy.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var products []Product
	for i, record := range records {
		if i == 0 {
			continue
		}

		soni, _ := strconv.Atoi(record[3])
		kelishNarxi := parseFloat(record[4])
		sotishNarxi := parseFloat(record[5])
		razmer, _ := strconv.Atoi(record[2])

		products = append(products, Product{
			TovarNomi:   record[1],
			Razmer:      razmer,
			Soni:        soni,
			KelishNarxi: kelishNarxi,
			SotishNarxi: sotishNarxi,
		})
	}

	for _, p := range products {
		fmt.Printf("%+v\n", p)
	}
}

func parseFloat(s string) float64 {
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "UZS", "")
	s = strings.TrimSpace(s)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

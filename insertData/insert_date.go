package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
)

type Product struct {
	TovarNomi   string
	Razmer      int
	Soni        int
	KelishNarxi float64
	SotishNarxi float64
}

func main() {
	// 📂 CSV ni ochamiz
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
			continue // header qatorni tashlab ketamiz
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

	// 🗄️ Bazaga ulanamiz
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:1111@db:5432/shop")
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	categoryID := "a5bd6076-717d-4d61-947c-c0a252e00abe"
	productType := "ml"

	// 📥 Insert qilish
	for _, p := range products {
		_, err := conn.Exec(context.Background(),
			`INSERT INTO products (name, size, type, price, selling_price, img_url, count, category_id, description)
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			p.TovarNomi,   // name
			p.Razmer,      // size
			productType,   // type
			p.KelishNarxi, // price
			p.SotishNarxi, // selling_price
			"",            // img_url (hozircha bo‘sh qoldirdim)
			p.Soni,        // count
			categoryID,    // category_id
			"",            // description (hozircha bo‘sh)
		)
		if err != nil {
			fmt.Println("Insert error:", err)
		}
	}

	fmt.Println("Insert tugadi ✅")
}

func parseFloat(s string) float64 {
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "UZS", "")
	s = strings.TrimSpace(s)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

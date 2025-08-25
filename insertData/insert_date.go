package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

func main() {
	// Excel fayl nomi
	excelFile := "AVTO DOKON TOVARLAR.xlsx"
	csvFile := "product.csv"

	// Excelni ochish
	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		log.Fatalf("Excel ochilmadi: %v", err)
	}

	// Sheet nomini olish (birinchi sheet)
	sheetName := f.GetSheetName(f.GetActiveSheetIndex())

	// Sheetdagi barcha satrlarni olish
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Satrlani o'qishda xato: %v", err)
	}

	// CSV fayl yaratish
	outFile, err := os.Create(csvFile)
	if err != nil {
		log.Fatalf("CSV fayl yaratishda xato: %v", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Har bir satrni CSV ga yozish
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			log.Fatalf("CSV ga yozishda xato: %v", err)
		}
	}

	fmt.Printf("✅ Muvaffaqiyatli CSV ga o‘tkazildi: %s\n", csvFile)
}

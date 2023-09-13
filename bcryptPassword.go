package main

import (
	"fmt"
	"log"
	"sync"
    "time"
	"github.com/tealeg/xlsx"
	"golang.org/x/crypto/bcrypt"
)

var mutex = &sync.Mutex{}

// Fungsi untuk mengenkripsi password
func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPassword)
}


func processRow(row *xlsx.Row, headerRow *xlsx.Row) {
	cells := row.Cells
	headerCells := headerRow.Cells

	// Iterasi melalui sel-sel dalam baris data
	for i, cell := range cells {
		columnName := headerCells[i].String() // Ambil nama kolom dari baris header

		// Periksa apakah nama kolom adalah "password"
		if columnName == "password" {
			// Hash nilai password
			value := cell.String()
			hashedPassword := hashPassword(value)
			mutex.Lock()
			cell.SetString(hashedPassword)
			mutex.Unlock()
		}
	}
}

func main() {
	startTime := time.Now()
	excelFileName := `C:\Users\LENOVO\Pictures\Camera Roll\6. request akun_12092023 v1.2.xlsx`

	// Baca file Excel
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		log.Fatal(err)
	}

	// Dapatkan sheet pertama dari file Excel
	sheet := xlFile.Sheets[0]

	// Ambil baris pertama (header)
	headerRow := sheet.Rows[0]

	// Buat worker pool
	workerCount := 4
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for _, row := range sheet.Rows[1:] { // Mulai dari baris kedua untuk menghindari header
				if row != nil {
					processRow(row, headerRow)
				}
			}
		}(i)
	}

	wg.Wait()

	// Tulis kembali ke file Excel
	outputFileName := `C:\Users\LENOVO\Pictures\Camera Roll\generate.xlsx`

	err = xlFile.Save(outputFileName)
	elapsedTime := time.Since(startTime)
	fmt.Printf("dalam waktu %v\n", elapsedTime)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File Excel berhasil diubah.")
}


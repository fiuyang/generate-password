package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/tealeg/xlsx"
	"golang.org/x/crypto/bcrypt"
)

// Worker pool
type Worker struct {
	ID        int
	JobQueue  chan *xlsx.Row
	ResultMap map[int]*xlsx.Row
	Wg        *sync.WaitGroup
}

// Fungsi untuk mengenkripsi password
func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPassword)
}

// Fungsi untuk memproses setiap baris pada file Excel
func processRow(row *xlsx.Row) *xlsx.Row {
	cells := row.Cells

	passwordCell := cells[1]
	password := passwordCell.String()

	hashedPassword := hashPassword(password)
	passwordCell.SetString(hashedPassword)

	return row
}

// Fungsi untuk menjalankan pekerjaan pemrosesan pada worker
func (w *Worker) Process() {
	for row := range w.JobQueue {
		processedRow := processRow(row)
		w.ResultMap[row.GetCoordinate()] = processedRow
		w.Wg.Done()
	}
}

// Fungsi utama
func main() {
	excelFileName := "Dealer.xlsx"

	// Baca file Excel
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		log.Fatal(err)
	}

	// Dapatkan sheet pertama dari file Excel
	sheet := xlFile.Sheets[0]

	// Buat job queue dan result map
	jobQueue := make(chan *xlsx.Row, len(sheet.Rows))
	resultMap := make(map[int]*xlsx.Row)

	// Buat worker pool
	workerCount := 4 // Jumlah worker yang diinginkan
	wg := &sync.WaitGroup{}
	wg.Add(len(sheet.Rows))
	workers := make([]*Worker, workerCount)

	for i := 0; i < workerCount; i++ {
		worker := &Worker{
			ID:        i,
			JobQueue:  jobQueue,
			ResultMap: resultMap,
			Wg:        wg,
		}
		workers[i] = worker
		go worker.Process()
	}

	// Tambahkan pekerjaan ke job queue
	for _, row := range sheet.Rows {
		jobQueue <- row
	}

	// Tutup job queue
	close(jobQueue)

	// Tunggu pemrosesan selesai
	wg.Wait()

	// Tulis kembali ke file Excel
	outputFileName := "Dealer.xlsx"
	for _, row := range resultMap {
		err = row.Sheet().Row(row.RowIndex).ReplaceRow(row)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = xlFile.Save(outputFileName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File Excel berhasil diubah.")
}

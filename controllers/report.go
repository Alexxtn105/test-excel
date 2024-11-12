package controllers

import (
	"bytes"
	"encoding/csv"
	"github.com/xuri/excelize/v2"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"test-excel/models"
	"time"
)

// PageData Define the data structure to pass to the template
type PageData struct {
	Title     string
	UserStats *[]models.UserStats
}

// ReportHandler handler function to respond to GET requests
func ReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Define the data to pass to the template
		data := PageData{
			Title:     "Report",
			UserStats: models.GetUserStats(),
		}

		// Parse and execute the template
		tmpl, err := template.ParseFiles("views/home.html")
		if err != nil {
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			log.Printf("Error loading template: %v", err)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Unable to render template", http.StatusInternalServerError)
			log.Printf("Error rendering template: %v", err)
		}
	} else {
		// Method not allowed
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ReportCSVHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write the header
	writer.Write([]string{"UserID", "UserName", "LoginCount", "LastLogin", "Active"})

	userStats := models.GetUserStats()
	// Write the data
	for _, user := range *userStats {
		writer.Write([]string{
			strconv.Itoa(user.UserID),
			user.UserName,
			strconv.Itoa(user.LoginCount),
			user.LastLogin.Format(time.RFC3339),
			strconv.FormatBool(user.Active),
		})
	}

	writer.Flush()

	// Set the response headers
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=user_stats.csv")
	w.Write(buf.Bytes())
}

func ReportCSVDownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		log.Printf("Error reading file: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		http.Error(w, "Unable to parse CSV", http.StatusBadRequest)
		log.Printf("Error parsing CSV: %v", err)
		return
	}

	err = models.AddRows(&rows)
	if err != nil {
		http.Error(w, "Failed to add rows", http.StatusInternalServerError)
		log.Printf("Failed to add rows: %v", err)
		return
	}

	http.Redirect(w, r, "/report", http.StatusSeeOther)
}

func ReportExcelHandler(w http.ResponseWriter, r *http.Request) {

	f := excelize.NewFile()

	// создаем новый лист, сохраняем индекс в переменной index
	index, _ := f.NewSheet("Sheet1")

	// пишем заголовки столбцов
	headers := []string{"User_ID", "User_Name", "LoginCount", "LastLogin", "Active"}
	for i, header := range headers {
		col := string(rune('A' + i)) // ASCII символа A == 65. Для i=0 col=A, i=1 col=B
		f.SetCellValue("Sheet1", col+"1", header)
	}

	// берем данные из БД
	userStats := models.GetUserStats()

	// пишем данные
	for i, user := range *userStats {
		row := strconv.Itoa(i + 2) // делаем +2 потосу что счет идет с 1 и первый столбец это заголовок
		f.SetCellInt("Sheet1", "A"+row, user.UserID)
		f.SetCellValue("Sheet1", "B"+row, user.UserName)
		f.SetCellInt("Sheet1", "C"+row, user.LoginCount)
		f.SetCellValue("Sheet1", "D"+row, user.LastLogin.Format("2006-01-02 15:04:05"))
		f.SetCellBool("Sheet1", "E"+row, user.Active)
	}

	// устанавливаем активный лист
	f.SetActiveSheet(index)

	// пишем файл в ответ
	// пишем заголовки
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment;filename=user_stats.xlsx")
	//пишем данные
	if err := f.Write(w); err != nil {
		http.Error(w, "Unable to write file", http.StatusInternalServerError)
		log.Printf("Error Writing file: %v", err)
	}
}

func ReportExcelDownloadHandler(w http.ResponseWriter, r *http.Request) {

}

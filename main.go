package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xuri/excelize/v2"
)

type SheetData struct {
	SheetName string              `json:"sheet_name"`
	Data      []map[string]string `json:"data"`
}

type ExcelResponse struct {
	Data []SheetData `json:"data"`
}

func loadExcelHandler(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := excelize.OpenFile("./data/" + filename)
		if err != nil {
			http.Error(w, "erro ao abrir arquivo", 500)
			return
		}

		var resp ExcelResponse

		for _, sheet := range f.GetSheetList() {
			rows, err := f.GetRows(sheet)
			if err != nil {
				continue
			}

			if len(rows) < 1 {
				continue
			}

			headers := rows[0]
			var sheetRows []map[string]string

			for _, row := range rows[1:] {
				item := map[string]string{}
				for i, cell := range row {
					if i < len(headers) {
						item[headers[i]] = cell
					}
				}
				sheetRows = append(sheetRows, item)
			}

			resp.Data = append(resp.Data, SheetData{
				SheetName: sheet,
				Data:      sheetRows,
			})
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func main() {
	fmt.Println("Backend rodando em: http://localhost:8080")

	http.HandleFunc("/data/Planilha1.xlsx", loadExcelHandler("Planilha1.xlsx"))
	http.HandleFunc("/data/Planilha2.xlsx", loadExcelHandler("Planilha2.xlsx"))
	http.HandleFunc("/data/Planilha3.xlsx", loadExcelHandler("Planilha3.xlsx"))

	http.ListenAndServe(":8080", nil)
}

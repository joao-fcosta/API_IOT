package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type Harmonic struct {
	Harmonic   float64 `json:"harmonic"`
	VMagnitude float64 `json:"V_magnitude"`
	IMagnitude float64 `json:"I_magnitude"`
}

func toFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func harmonicsHandler(w http.ResponseWriter, r *http.Request) {
	f, err := excelize.OpenFile("./data/Planilha2.xlsx")
	if err != nil {
		http.Error(w, "erro ao abrir planilha", 500)
		return
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		http.Error(w, "erro ao ler planilha", 500)
		return
	}

	var data []Harmonic
	for _, row := range rows[1:] {
		if len(row) < 3 {
			continue
		}
		d := Harmonic{
			Harmonic:   toFloat(row[0]),
			VMagnitude: toFloat(row[1]),
			IMagnitude: toFloat(row[2]),
		}
		data = append(data, d)
	}

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/harmonics", harmonicsHandler)
	http.ListenAndServe(":8080", nil)
}

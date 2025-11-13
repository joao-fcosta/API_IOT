package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

// Constantes
const pastaBase = "data"

// Estrutura de dados para a guia
type SheetData struct {
	SheetName string `json:"sheet_name"`
	// O campo Data agora é um slice de mapas (objetos),
	// onde a chave (string) é o nome da coluna (cabeçalho)
	Data []map[string]string `json:"data"`
}

// Estrutura para a resposta da API (reutilizada)
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    []SheetData `json:"data"`
}

// -------------------------------------------------------------
// FUNÇÃO CENTRAL: LÊ O ARQUIVO EXCEL E MAPEIA DADOS
// --------------------------------------------------------------------------------

// readExcelData lê o arquivo e retorna os dados da primeira guia mapeados como objetos JSON.
func readExcelData(filePath string) ([]SheetData, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o arquivo '%s': %w", filePath, err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("o arquivo Excel não contém nenhuma guia")
	}

	// Obtém todas as linhas da guia
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler as linhas da guia '%s': %w", sheetName, err)
	}

	if len(rows) < 1 {
		return nil, fmt.Errorf("a guia '%s' está vazia ou não tem cabeçalhos", sheetName)
	}

	// 1. Identifica os Cabeçalhos (primeira linha)
	headers := rows[0]
	// 2. Separa os Dados (da segunda linha em diante)
	dataRows := rows[1:]

	// Slice para armazenar os dados mapeados (Array de Objetos)
	var mappedData []map[string]string

	// 3. Mapeamento de Linha para Objeto (Chave/Valor)
	for _, row := range dataRows {
		record := make(map[string]string)
		for i, header := range headers {
			if i < len(row) {
				// Mapeia o cabeçalho (header) para o valor da célula (row[i])
				record[header] = row[i]
			} else {
				// Trata caso a linha seja mais curta que os cabeçalhos
				record[header] = ""
			}
		}
		mappedData = append(mappedData, record)
	}

	// 4. Retorna o resultado final
	result := []SheetData{
		{
			SheetName: sheetName,
			Data:      mappedData, // Dados no novo formato de objeto
		},
	}

	return result, nil
}

// -------------------------------------------------------------
// HANDLER HTTP
// -------------------------------------------------------------

func getFileDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Método não permitido. Use GET.", http.StatusMethodNotAllowed)
		return
	}

	filePath := r.URL.Path
	filename := strings.TrimPrefix(filePath, "/data/")

	if filename == "" || filename == "/" {
		http.Error(w, "Nome do arquivo não especificado. Exemplo: /data/Planilha1.xlsx", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(pastaBase, filename)

	data, err := readExcelData(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, fmt.Sprintf("Arquivo não encontrado em: %s", fullPath), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Erro ao processar o arquivo: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}

	response := APIResponse{
		Status:  "success",
		Message: fmt.Sprintf("Dados extraídos com sucesso da primeira guia de '%s'", filename),
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/data/", getFileDataHandler)

	port := ":8080"
	fmt.Printf("🚀 Servidor iniciado na porta %s. Para acessar, use: http://localhost%s/data/{nome-do-arquivo.xlsx}\n", port, port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Erro ao iniciar o servidor: %s\n", err)
		os.Exit(1)
	}
}

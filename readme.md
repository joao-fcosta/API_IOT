# API Go para Leitura de Excel Local 📍
Esta documentação fornece os passos essenciais para rodar a API em Go que lê um arquivo Excel (`.xlsx`) da pasta `data/` e expõe seus dados via HTTP.
Nesse caso, realizaremos a leitura dos dados obtidos pelo circuito que captura informações da corrente elétrica em tempo real. ⚡

> [!NOTE]
> 🛠️ Pré-requisitos - Para executar esta aplicação, você precisa ter instalado: <br>
> **Go Language** Versão **1.16** ou superior. <br>
> **Postman** ou `curl` (Para testar o endpoint da API).

---

## 🚀 1. Como Rodar a Aplicação
Execute o servidor a partir do diretório raiz do projeto:

```bash
go run main.go
```

O servidor será iniciado na porta **8080**.


## 🧪 2. Como Acessar os Dados (Testes)

A API usa uma rota **GET** para ler o arquivo e retorna os dados da **primeira guia** em formato JSON (chave/valor).

#### Endpoint de Acesso

Você deve passar o nome do arquivo na URL. Use o seu nome de arquivo.

| Método | URL (Exemplo) |
| :--- | :--- |
| **GET** | `http://localhost:8080/data/Planilha1.xlsx` |

#### A. Usando `curl`

Para testar o arquivo no terminal:

```bash
curl -X GET http://localhost:8080/data/Planilha1.xlsx
```

#### B. Usando Postman

1.  Crie uma nova requisição.
2.  Defina o **Método** como **`GET`**.
3.  Cole a URL com o nome do seu arquivo.
4.  Clique em **`Send`**.

O retorno será o JSON com os dados da primeira guia do arquivo Excel.

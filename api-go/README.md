# FitTracker API (Go Migration)

Este é o backend do FitTracker migrado para Go utilizando a **Standard Library** o máximo possível.

## Requisitos
- Go 1.22 ou superior
- Ollama (opcional, para IA)

## Como rodar
1. Entre no diretório: `cd api-go`
2. Rode o servidor: `go run ./cmd/api/main.go`

O servidor iniciará na porta `:5000`.

## Diferenças da Versão C#
- **Performance:** Binário nativo e leve.
- **Simplicidade:** Sem frameworks pesados, usando o `http.ServeMux` nativo.
- **Banco de Dados:** SQLite via `database/sql` (driver puro Go).
- **IA:** Chamadas diretas para a API do Ollama.

## Estrutura do Projeto
- `cmd/api/`: Ponto de entrada e rotas.
- `internal/handlers/`: Lógica dos endpoints (Slices).
- `internal/models/`: Estruturas de dados e JSON.
- `internal/auth/`: JWT e Hashing de senhas.
- `internal/database/`: Conexão e Schema SQL.
- `internal/ai/`: Integração com Ollama.

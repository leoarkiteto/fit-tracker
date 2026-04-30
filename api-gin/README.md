# FitTracker API (Go + Gin + GORM)

Este é o backend do FitTracker migrado para Go utilizando o framework **Gin** e o ORM **GORM**.

## Requisitos
- Go 1.22 ou superior
- Ollama (opcional, para IA)

## Como rodar
1. Entre no diretório: `cd api-gin`
2. Rode o servidor: `go run ./cmd/api/main.go`

O servidor iniciará na porta `:5001` (para não conflitar com a versão Standard Library na `:5000`).

## Diferenças da Versão Standard (api-go)
- **Framework:** Usa Gin para roteamento, middleware e validação de JSON.
- **ORM:** Usa GORM para manipulação do banco de dados, facilitando relacionamentos (Exercises dentro de Workouts).
- **Auto-Migration:** O schema do banco de dados é gerado automaticamente a partir das structs Go.
- **Produtividade:** Código mais conciso para operações CRUD e tratamento de parâmetros de rota.

## Estrutura do Projeto
- `cmd/api/`: Ponto de entrada e configuração de rotas.
- `internal/handlers/`: Controladores Gin (recebem *gin.Context).
- `internal/models/`: Estruturas de dados com tags GORM.
- `internal/auth/`: Lógica de JWT e Criptografia.
- `internal/database/`: Conexão e AutoMigrate.
- `internal/ai/`: Integração com Ollama via LangChainGo.

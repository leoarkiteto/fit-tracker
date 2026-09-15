# FitTracker API (Go Migration)

Este é o backend do FitTracker migrado para Go utilizando a **Standard Library** o máximo possível.

## Requisitos
- Go 1.22 ou superior
- Ollama (opcional, para IA)
- Variável de ambiente `JWT_SECRET` (obrigatória: o servidor não inicia sem ela)

## Como rodar
1. Entre no diretório: `cd backend`
2. Exporte o segredo: `export JWT_SECRET=my-secret`
3. Rode o servidor: `go run ./cmd/api`

O servidor iniciará na porta `:5000`.

## Testes
`go test ./...` roda os testes de caso de uso de cada slice. Eles não usam HTTP nem banco: cada
teste injeta os dublês em memória de `internal/<slice>/mock`.

## Arquitetura

O backend segue **vertical slice architecture**: cada feature é um pacote em `internal/` que reúne
todas as suas camadas, do domínio ao acesso a dados.

```
backend/
├── cmd/api/main.go        # composition root: monta banco, signer e routers
├── pkg/                   # infraestrutura compartilhada
│   ├── apperror/          # erro de domínio + status HTTP
│   ├── auth/              # bcrypt + JWT
│   ├── database/          # conexão SQLite e schema
│   └── httpx/             # guarda de autenticação, body limit, JSON, CORS
└── internal/
    ├── user/              # registro, login, /me, troca de senha
    ├── profile/           # perfil do usuário
    ├── workout/           # treinos e exercícios
    ├── completedworkout/  # histórico e estatísticas
    ├── bioimpedance/      # bioimpedância
    ├── water/             # consumo de água
    └── aiplanning/        # planejamento com Ollama
```

Cada slice tem sempre o mesmo formato:

| Arquivo | Papel |
|---------|-------|
| `<slice>.go` | entidade + erros de domínio |
| `port.go` | interfaces de repositório/serviços externos (ports) |
| `service/` | casos de uso — um arquivo por caso de uso |
| `handler/` | HTTP — `handler.go` com `Router(mux, signer, svc)` e um arquivo por caso de uso |
| `infrastructure/` | implementação SQLite dos ports |
| `mock/` | dublês em memória usados pelos testes |

O `service` depende apenas do domínio e dos ports, então é testável sem HTTP e sem banco.

## Diferenças da Versão C#
- **Performance:** Binário nativo e leve.
- **Simplicidade:** Sem frameworks pesados, usando o `http.ServeMux` nativo.
- **Banco de Dados:** SQLite via `database/sql` (driver puro Go).
- **IA:** Chamadas diretas para a API do Ollama.

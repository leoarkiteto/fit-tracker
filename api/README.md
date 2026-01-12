# FitTracker API

API RESTful para o aplicativo FitTracker, desenvolvida com ASP.NET Core Minimal API, Entity Framework Core e SQLite.

## Requisitos

- .NET 9.0 SDK ou superior

## Executando a API

```bash
cd FitTracker.Api
dotnet run
```

A API estará disponível em:
- **HTTP**: http://localhost:5000
- **HTTPS**: https://localhost:5001

O Swagger UI estará disponível na raiz: http://localhost:5000/

## Estrutura do Projeto

```
FitTracker.Api/
├── Models/           # Entidades do banco de dados
├── Data/             # DbContext e configurações EF Core
├── DTOs/             # Data Transfer Objects
├── Endpoints/        # Minimal API endpoints
├── Program.cs        # Ponto de entrada e configuração
└── appsettings.json  # Configurações da aplicação
```

## Endpoints Disponíveis

### Profiles (Perfis de Usuário)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/profiles` | Lista todos os perfis |
| GET | `/api/profiles/{id}` | Obtém um perfil por ID |
| POST | `/api/profiles` | Cria um novo perfil |
| PUT | `/api/profiles/{id}` | Atualiza um perfil |
| DELETE | `/api/profiles/{id}` | Remove um perfil |

### Workouts (Treinos)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/profiles/{profileId}/workouts` | Lista todos os treinos |
| GET | `/api/profiles/{profileId}/workouts/{id}` | Obtém um treino por ID |
| GET | `/api/profiles/{profileId}/workouts/today` | Obtém treinos de hoje |
| POST | `/api/profiles/{profileId}/workouts` | Cria um novo treino |
| PUT | `/api/profiles/{profileId}/workouts/{id}` | Atualiza um treino |
| DELETE | `/api/profiles/{profileId}/workouts/{id}` | Remove um treino |

### Bioimpedance (Bioimpedância)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/profiles/{profileId}/bioimpedance` | Lista todo o histórico |
| GET | `/api/profiles/{profileId}/bioimpedance/{id}` | Obtém um registro por ID |
| GET | `/api/profiles/{profileId}/bioimpedance/latest` | Obtém o registro mais recente |
| POST | `/api/profiles/{profileId}/bioimpedance` | Cria um novo registro |
| PUT | `/api/profiles/{profileId}/bioimpedance/{id}` | Atualiza um registro |
| DELETE | `/api/profiles/{profileId}/bioimpedance/{id}` | Remove um registro |

### Completed Workouts (Treinos Concluídos)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/profiles/{profileId}/completed-workouts` | Lista todos os treinos concluídos |
| GET | `/api/profiles/{profileId}/completed-workouts/stats` | Obtém estatísticas |
| POST | `/api/profiles/{profileId}/completed-workouts` | Marca um treino como concluído |
| DELETE | `/api/profiles/{profileId}/completed-workouts/{id}` | Remove um registro |

### Health Check

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/health` | Verifica status da API |

## Exemplos de Uso

### Criar um Perfil

```bash
curl -X POST http://localhost:5000/api/profiles \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "age": 30,
    "height": 175,
    "currentWeight": 80,
    "goalWeight": 75
  }'
```

### Criar um Treino

```bash
curl -X POST http://localhost:5000/api/profiles/{profileId}/workouts \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Treino de Peito",
    "description": "Treino focado em hipertrofia",
    "goal": "hypertrophy",
    "days": ["monday", "thursday"],
    "exercises": [
      {
        "name": "Supino Reto",
        "muscleGroup": "chest",
        "sets": 4,
        "reps": 12,
        "weight": 60,
        "restSeconds": 90
      }
    ]
  }'
```

### Registrar Bioimpedância

```bash
curl -X POST http://localhost:5000/api/profiles/{profileId}/bioimpedance \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2026-01-12T10:00:00Z",
    "weight": 80.5,
    "bodyFatPercentage": 18.5,
    "muscleMass": 35.2,
    "boneMass": 3.5,
    "waterPercentage": 55.0,
    "visceralFat": 8,
    "bmr": 1800,
    "metabolicAge": 28
  }'
```

## Banco de Dados

O banco de dados SQLite (`fittracker.db`) é criado automaticamente na primeira execução. Não é necessário rodar migrations manualmente.

## Configuração CORS

A API está configurada para aceitar requisições de qualquer origem (`AllowAll`), facilitando o desenvolvimento com o app mobile.

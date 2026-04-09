# BudgetGen — Backend (Go)

## Stack
- **Go 1.22**
- **Gin** — HTTP framework
- **GORM** — ORM
- **PostgreSQL** — banco de dados
- **JWT** — autenticação

## Setup

### 1. Instalar Go
Baixe em https://go.dev/dl/ e instale o `.msi` para Windows.

### 2. Configurar variáveis de ambiente
```bash
cp .env.example .env
# edite o .env com suas credenciais
```

### 3. Instalar dependências
```bash
go mod tidy
```

### 4. Rodar
```bash
go run ./cmd/api
```

## Endpoints

### Auth (público)
| Método | Rota | Descrição |
|--------|------|-----------|
| POST | `/auth/register` | Cadastro |
| POST | `/auth/login` | Login |

### API (requer `Authorization: Bearer <token>`)
| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/me` | Dados do usuário |
| GET/POST | `/api/quotes` | Listar / Criar orçamentos |
| GET/PUT/DELETE | `/api/quotes/:id` | Orçamento específico |
| GET/POST | `/api/templates` | Listar / Criar templates |
| GET/PUT/DELETE | `/api/templates/:id` | Template específico |
| GET/PUT | `/api/settings` | Configurações da empresa |

# TaskFlow

Gerenciador de projetos com quadro Kanban, sistema de tarefas com prioridade e prazo, e dashboard pessoal de tarefas atribuídas.

**Stack:** Go 1.25+ · Gin · GORM · PostgreSQL · Bootstrap 5

---

## Opção A — Com Docker (recomendado)

**Pré-requisito:** [Docker Desktop](https://www.docker.com/products/docker-desktop/)

```bash
# 1. Clone o repositório
git clone https://github.com/ThiagoGama1/Task_Flow
cd Task_Flow

# 2. Crie o arquivo de ambiente
Copy-Item .env.example .env      # Windows (PowerShell)
# cp .env.example .env           # Linux / macOS

# 3. Suba banco + aplicação
docker compose up -d
```

Acesse: **http://localhost:3000**

As tabelas são criadas automaticamente. Para parar: `docker compose down`

### Popular com dados de exemplo

```bash
docker compose exec app ./seed
```

Cria 3 usuários, 2 projetos, 12 tarefas e comentários prontos para demonstração.

| E-mail                  | Senha   |
|-------------------------|---------|
| demo@taskflow.app       | demo123 |
| ana@taskflow.app        | demo123 |
| carlos@taskflow.app     | demo123 |

---

## Opção B — Sem Docker (Go + PostgreSQL local)

**Pré-requisitos:** [Go 1.25+](https://go.dev/dl/) e [PostgreSQL](https://www.postgresql.org/download/)

```bash
# 1. Clone o repositório
git clone https://github.com/ThiagoGama1/Task_Flow
cd Task_Flow

# 2. Crie o banco de dados
psql -U postgres -c "CREATE USER taskflow WITH PASSWORD 'taskflow';"
psql -U postgres -c "CREATE DATABASE taskflow OWNER taskflow;"

# 3. Crie o arquivo de ambiente
Copy-Item .env.example .env      # Windows (PowerShell)
# cp .env.example .env           # Linux / macOS

# 4. Instale as dependências e rode
go mod tidy
go run ./cmd/server
```

Acesse: **http://localhost:3000**

### Popular com dados de exemplo

```bash
go run ./cmd/seed
```

---

## Testes

```bash
go test ./tests/... -v
```

---

## Estrutura do projeto

```
taskflow/
├── cmd/
│   ├── server/         # Ponto de entrada da aplicação
│   └── seed/           # Populador de dados de exemplo
├── internal/
│   ├── app/            # Inicialização (rotas, templates, middlewares)
│   ├── config/         # Leitura de variáveis de ambiente
│   ├── database/       # Conexão e migrations automáticas
│   ├── handlers/       # Controllers HTTP (auth, projects, tasks, dashboard)
│   ├── middleware/      # Autenticação de sessão
│   ├── models/         # Entidades GORM (User, Project, Task, Comment, ActivityLog)
│   ├── repositories/   # Camada de acesso ao banco de dados
│   └── routes/         # Registro de rotas por domínio
├── static/css/         # CSS customizado
├── templates/          # Templates HTML (Go html/template)
├── tests/              # Testes de integração
├── .env.example        # Variáveis de ambiente necessárias
└── go.mod
```

---

## Funcionalidades

- **Autenticação** — cadastro e login com senha criptografada (bcrypt)
- **Projetos** — criação, membros por convite, exclusão
- **Tarefas** — prioridade (baixa/média/alta), prazo, responsável, status
- **Kanban** — colunas A Fazer / Em Andamento / Concluído com drag-and-drop e ordenação
- **Dashboard** — tarefas atribuídas agrupadas por urgência (atrasadas, vencem hoje, próximas)

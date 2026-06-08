# TaskFlow

Gerenciador de projetos com quadro Kanban, sistema de tarefas com prioridade e prazo, e dashboard pessoal de tarefas atribuídas.

**Stack:** Go 1.21+ · Gin · GORM · PostgreSQL · Bootstrap 5

---

## Pré-requisitos

- **Go 1.21+** — https://go.dev/dl/
- **PostgreSQL** — https://www.postgresql.org/download/

Para verificar se já estão instalados:

```bash
go version
psql --version
```

---

## Como rodar

### 1. Clone o repositório

```bash
git clone <url-do-repositorio>
cd taskflow
```

### 2. Crie o banco de dados no PostgreSQL

```sql
-- No psql ou pgAdmin:
CREATE USER taskflow WITH PASSWORD 'taskflow';
CREATE DATABASE taskflow OWNER taskflow;
```

Ou em uma linha pelo terminal:

```bash
psql -U postgres -c "CREATE USER taskflow WITH PASSWORD 'taskflow';"
psql -U postgres -c "CREATE DATABASE taskflow OWNER taskflow;"
```

### 3. Configure o ambiente

Copie o arquivo de exemplo:

```bash
# Linux / macOS
cp .env.example .env

# Windows (PowerShell)
Copy-Item .env.example .env
```

Edite o `.env` com as credenciais do seu banco:

```
DATABASE_URL=postgres://taskflow:taskflow@localhost:5432/taskflow?sslmode=disable
SESSION_SECRET=taskflow-dev-secret-2024
PORT=3000
GIN_MODE=debug
```

### 4. Baixe as dependências

```bash
go mod tidy
```

### 5. Rode o servidor
cd "C:\Users\thiagogama\Desktop\des web\taskflow"

```bash
go run ./cmd/server
```

Acesse em: **http://localhost:3000**

As tabelas são criadas automaticamente na primeira execução (GORM AutoMigrate). Não é necessário rodar migrations manualmente.

---

## Estrutura do projeto

```
taskflow/
├── cmd/server/         # Ponto de entrada da aplicação
├── internal/
│   ├── app/            # Inicialização (rotas, templates, middlewares)
│   ├── config/         # Leitura de variáveis de ambiente
│   ├── database/       # Conexão e migrations automáticas
│   ├── handlers/       # Controllers HTTP (auth, projects, tasks, dashboard)
│   ├── middleware/      # Autenticação de sessão
│   ├── models/         # Entidades GORM (User, Project, Task)
│   ├── repositories/   # Camada de acesso ao banco de dados
│   └── routes/         # Registro de rotas por domínio
├── static/css/         # CSS customizado
├── templates/          # Templates HTML (Go html/template)
│   ├── auth/
│   ├── layout/
│   ├── projects/
│   └── tasks/
├── tests/              # Testes de integração
├── .env.example        # Variáveis de ambiente necessárias
└── go.mod
```

---

## Funcionalidades

- **Autenticação** — cadastro e login com senha criptografada (bcrypt)
- **Projetos** — criação, membros por convite, exclusão
- **Tarefas** — prioridade (baixa/média/alta), prazo, responsável, status
- **Kanban** — visualização em colunas A Fazer / Em Andamento / Concluído
- **Dashboard** — tarefas atribuídas ao usuário agrupadas por urgência

---

## Rodando os testes

```bash
go test ./tests/... -v
```

---

## Compilando um executável

Para gerar um binário único que pode ser executado sem o Go instalado:

```bash
go build -o taskflow ./cmd/server

# Windows
go build -o taskflow.exe ./cmd/server
```

Execute com:

```bash
./taskflow        # Linux / macOS
./taskflow.exe    # Windows
```

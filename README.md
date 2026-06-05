# 🏭 Planeta Azul — Sistema de Gestão Industrial

Monorepo com backend Go + Gin e frontend Next.js 14.

```
planeta-azul/
├── backend/          Go + Gin  (porta 8080)
│   ├── cmd/api/      main.go
│   ├── internal/
│   │   ├── auth/     JWT RS256
│   │   ├── handlers/ auth, items, users
│   │   ├── middleware/ RequireAuth, RequireRoles
│   │   ├── models/   todos os tipos + DTOs
│   │   ├── repository/ MemStore (sem DB) / futuro: PostgreSQL
│   │   └── config/
│   └── migrations/   001_init.sql
└── frontend/         Next.js 14 + Tailwind + shadcn/ui  (porta 3000)
    └── src/app/
        ├── login/        Tela de login
        ├── dashboard/    Home + layout com sidebar
        ├── tasks/        Módulo 1 — lista, detalhe, novo item
        ├── access/       Módulo 2 (placeholder)
        ├── data/         Módulo 3 (placeholder)
        ├── forms/        Módulo 4 (placeholder)
        └── ocr/          Módulo 5 (placeholder)
```

---

## Rodando sem banco de dados (modo dev)

### Backend

```bash
cd backend

# Instalar dependências
go mod tidy

# Rodar (gera chaves JWT efêmeras automaticamente)
go run ./cmd/api
# → http://localhost:8080
# → GET http://localhost:8080/health
```

Credenciais de seed:
| Email                      | Senha     | Role        |
|----------------------------|-----------|-------------|
| admin@planetaazul.com      | admin123  | admin       |
| chefe@planetaazul.com      | chefe123  | chefe_area  |
| sup@planetaazul.com        | sup123    | supervisor  |
| membro@planetaazul.com     | membro123 | membro      |

### Frontend

```bash
cd frontend

npm install
npm run dev
# → http://localhost:3000
```

O Next.js faz proxy de `/api/*` e `/auth/*` para `localhost:8080` via `next.config.js`.

---

## Com banco de dados (futuro)

1. Criar banco PostgreSQL e rodar `backend/migrations/001_init.sql`
2. Copiar `.env.example` → `.env` e preencher variáveis
3. Gerar chaves RSA:
   ```bash
   openssl genrsa -out backend/scripts/private.pem 2048
   openssl rsa -in backend/scripts/private.pem -pubout -out backend/scripts/public.pem
   ```
4. Trocar `MemStore` por repositório PostgreSQL (GORM)

---

## Endpoints disponíveis

```
POST   /auth/login
POST   /auth/refresh
POST   /auth/logout

GET    /api/v1/auth/me
GET    /api/v1/users
GET    /api/v1/users/:id
GET    /api/v1/areas
GET    /api/v1/notifications

GET    /api/v1/items
POST   /api/v1/items
GET    /api/v1/items/:id
PATCH  /api/v1/items/:id
DELETE /api/v1/items/:id          (chefe_area+)
GET    /api/v1/items/:id/comments
POST   /api/v1/items/:id/comments
GET    /api/v1/items/:id/interactions
POST   /api/v1/items/:id/interactions
```

---

## Próximos passos

- [ ] Repositório PostgreSQL com GORM (substituir MemStore)
- [ ] Módulo 2 — gestão de acesso e sync com banco da fábrica
- [ ] Módulo 3 — KPIs via views materializadas
- [ ] Módulo 4 — form builder com JSONB
- [ ] Módulo 5 — proxy OCR
- [ ] SSE para notificações em tempo real
- [ ] Testes unitários (Go) e E2E (Playwright)

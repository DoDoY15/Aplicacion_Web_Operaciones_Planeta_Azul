# 🏭 Planeta Azul — Sistema de Gestão Industrial

Monorepo com backend Go + Gin e frontend Next.js 15 + React 19. Deploy em produção: backend no **Railway**, frontend no **Cloudflare Workers** (via adaptador OpenNext).

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
│   ├── migrations/   001_init.sql
│   └── railway.toml  Deploy no Railway
└── frontend/         Next.js 15 + React 19 + Tailwind + shadcn/ui  (porta 3000)
    ├── wrangler.jsonc Deploy no Cloudflare (Workers + OpenNext)
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

Em desenvolvimento, o Next.js faz proxy de `/api/*` e `/auth/*` para `localhost:8080` via `next.config.js`.

---

## Deploy em produção

### Backend — Railway

Configurado via [`backend/railway.toml`](backend/railway.toml) (builder Nixpacks, `go run ./cmd/api`, healthcheck em `/health`).

Variáveis de ambiente principais (ver [`backend/.env.example`](backend/.env.example)):
- `ALLOWED_ORIGINS` — deve incluir o domínio do frontend publicado no Cloudflare Pages (CORS via `gin-contrib/cors`)
- `JWT_*`, `DB_*`, `FABRICA_DB_*`, `OCR_SERVICE_URL`

### Frontend — Cloudflare Workers (via OpenNext)

Deploy feito com [`@opennextjs/cloudflare`](https://opennext.js.org/cloudflare), o adaptador
oficial mantido pela Cloudflare/Vercel — sucessor do depreciado `@cloudflare/next-on-pages`
(que só suportava o runtime "Edge" e parou de ser mantido). Com o OpenNext, o app roda em
runtime **Node.js completo** nos Cloudflare Workers, sem precisar de `export const runtime = 'edge'`.

```bash
cd frontend

# Build + deploy direto (gera .open-next/ e publica via wrangler)
npm run deploy

# Ou só para testar localmente, simulando o runtime real dos Workers:
npm run preview
```

Configuração em [`frontend/wrangler.jsonc`](frontend/wrangler.jsonc) (`main: ".open-next/worker.js"`,
`compatibility_flags: ["nodejs_compat"]`, assets servidos via binding `ASSETS`) e
[`frontend/open-next.config.ts`](frontend/open-next.config.ts) (config mínima — o projeto
não usa cache incremental via R2 nem outros bindings).

Variável de ambiente necessária no Cloudflare (configurar em **Workers & Pages → planeta-azul
→ Settings → Variables**, ver [`frontend/.env.example`](frontend/.env.example)):
- `NEXT_PUBLIC_API_URL` — URL pública do backend no Railway (sem barra final). Em dev local deixe vazio — o proxy do Next.js cuida disso.

> **Importante**: ao publicar uma alteração, sempre dispare um deploy **novo** (push para
> `main`, que aciona a integração com o Git, ou `npm run deploy`). Nunca use "Retry deployment"
> de um deploy antigo no painel — isso reaproveita o commit antigo daquele deploy em vez do
> HEAD atual da branch, e pode reintroduzir bugs já corrigidos.

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

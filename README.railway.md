# Deploy no Railway

## Pré-requisitos
- Conta no [Railway](https://railway.app)
- Repositório no GitHub com este código

---

## Passo a passo

### 1. Criar projeto no Railway
- Acesse railway.app → New Project → Deploy from GitHub repo
- Selecione o repositório `BackEnd-BudgetGenerator`

---

### 2. Adicionar PostgreSQL
- No projeto → Add Service → Database → PostgreSQL
- O Railway cria automaticamente a variável `DATABASE_URL`

---

### 3. Adicionar MinIO
- No projeto → Add Service → GitHub Repo → selecione o repo
- Em Settings → Build:
  - **Dockerfile Path**: `Dockerfile.minio`
- Em Settings → Networking:
  - Gerar domínio público na porta `9000` (API — para o frontend acessar as imagens)
- Em Settings → Volumes:
  - Adicionar volume em `/data` para persistir os arquivos

> ⚠️ O Railway não persiste dados por padrão — sem volume os arquivos somem ao reiniciar.

---

### 4. Serviço AUTH
- Add Service → GitHub Repo → selecione o repo
- Em Settings → Build:
  - **Dockerfile Path**: `Dockerfile.auth`
- Variáveis de ambiente:
  ```
  DATABASE_URL=${{Postgres.DATABASE_URL}}
  JWT_SECRET=<gere uma string aleatória longa>
  ```
- Em Settings → Networking → gerar domínio público

---

### 5. Serviço CORE
- Add Service → GitHub Repo → selecione o repo (segunda instância)
- Em Settings → Build:
  - **Dockerfile Path**: `Dockerfile.core`
- Variáveis de ambiente:
  ```
  DATABASE_URL=${{Postgres.DATABASE_URL}}
  JWT_SECRET=<mesma string do auth>
  MINIO_ENDPOINT=<domínio privado do MinIO>:9000
  MINIO_ACCESS_KEY=budgetgen
  MINIO_SECRET_KEY=budgetgen123
  MINIO_BUCKET=uploads
  MINIO_PUBLIC_URL=https://<domínio público do MinIO>
  ```
- Em Settings → Networking → gerar domínio público

---

### 6. Configurar o Frontend
Após os deploys, pegue as URLs públicas dos serviços e configure no frontend:

```env
VITE_AUTH_URL=https://<url-do-auth-service>
VITE_CORE_URL=https://<url-do-core-service>
```

---

## Variáveis obrigatórias por serviço

| Serviço | Variável | Valor |
|---------|----------|-------|
| auth | DATABASE_URL | referência ao Postgres |
| auth | JWT_SECRET | string aleatória |
| core | DATABASE_URL | referência ao Postgres |
| core | JWT_SECRET | mesma do auth |
| core | MINIO_ENDPOINT | host:porta do MinIO |
| core | MINIO_ACCESS_KEY | usuário MinIO |
| core | MINIO_SECRET_KEY | senha MinIO |
| core | MINIO_BUCKET | uploads |
| core | MINIO_PUBLIC_URL | URL pública do MinIO |

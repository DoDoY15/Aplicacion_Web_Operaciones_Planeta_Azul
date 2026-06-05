-- Planeta Azul — Migration 001: Schema inicial
-- Executar em: PostgreSQL do App (banco separado da fábrica)

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Enums
CREATE TYPE user_role AS ENUM ('admin', 'chefe_geral', 'chefe_area', 'supervisor', 'membro');
CREATE TYPE item_status AS ENUM ('draft', 'pending', 'in_progress', 'waiting', 'done', 'rejected');
CREATE TYPE item_visibility AS ENUM ('private', 'team', 'public');
CREATE TYPE item_priority AS ENUM ('low', 'medium', 'high', 'urgent');
CREATE TYPE approval_decision AS ENUM ('approved', 'rejected');
CREATE TYPE interaction_status AS ENUM ('open', 'resolved');

-- Areas
CREATE TABLE areas (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Users
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          VARCHAR(150) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role          user_role NOT NULL DEFAULT 'membro',
    cargo         VARCHAR(100),
    area_id       UUID REFERENCES areas(id) ON DELETE SET NULL,
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email   ON users(email);
CREATE INDEX idx_users_area_id ON users(area_id);

-- Items (self-referencing, entity única para Projeto/Task/Subtask)
CREATE TABLE items (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id         UUID REFERENCES items(id) ON DELETE CASCADE,
    title             VARCHAR(255) NOT NULL,
    description       TEXT,
    created_by        UUID NOT NULL REFERENCES users(id),
    assigned_to       UUID REFERENCES users(id),
    area_id           UUID NOT NULL REFERENCES areas(id),
    deleted_by        UUID REFERENCES users(id),
    status            item_status NOT NULL DEFAULT 'draft',
    visibility        item_visibility NOT NULL DEFAULT 'team',
    requires_approval BOOLEAN NOT NULL DEFAULT FALSE,
    priority          item_priority NOT NULL DEFAULT 'medium',
    deadline          TIMESTAMPTZ,
    completed_at      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ  -- soft delete
);

CREATE INDEX idx_items_area_id    ON items(area_id);
CREATE INDEX idx_items_created_by ON items(created_by);
CREATE INDEX idx_items_assigned_to ON items(assigned_to);
CREATE INDEX idx_items_parent_id  ON items(parent_id);
CREATE INDEX idx_items_status     ON items(status);
CREATE INDEX idx_items_deleted_at ON items(deleted_at);

-- Item assignments (multiple responsibles per item)
CREATE TABLE item_assignments (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_id      UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    user_id      UUID NOT NULL REFERENCES users(id),
    role_in_item VARCHAR(50) NOT NULL DEFAULT 'executor',
    assigned_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_by  UUID NOT NULL REFERENCES users(id),
    UNIQUE(item_id, user_id)
);

-- Interactions (devolução/aguardando decisão)
CREATE TABLE interactions (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_id      UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    opened_by    UUID NOT NULL REFERENCES users(id),
    addressed_to UUID NOT NULL REFERENCES users(id),
    message      TEXT NOT NULL,
    response     TEXT,
    status       interaction_status NOT NULL DEFAULT 'open',
    resolved_at  TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_interactions_item_id      ON interactions(item_id);
CREATE INDEX idx_interactions_addressed_to ON interactions(addressed_to);

-- Comments
CREATE TABLE comments (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_id    UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id),
    content    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_item_id ON comments(item_id);

-- Approvals
CREATE TABLE approvals (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_id    UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    reviewer   UUID NOT NULL REFERENCES users(id),
    decision   approval_decision NOT NULL,
    note       TEXT,
    decided_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_approvals_item_id ON approvals(item_id);

-- Area access (cruzamento de áreas)
CREATE TABLE area_access (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    granted_by UUID NOT NULL REFERENCES users(id),
    granted_to UUID NOT NULL REFERENCES users(id),
    area_id    UUID NOT NULL REFERENCES areas(id),
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(granted_to, area_id)
);

-- Notifications
CREATE TABLE notifications (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type       VARCHAR(50) NOT NULL,
    ref_id     UUID,
    ref_type   VARCHAR(50),
    message    TEXT NOT NULL,
    read       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_read    ON notifications(read);

-- Auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN NEW.updated_at = NOW(); RETURN NEW; END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_items_updated_at
BEFORE UPDATE ON items
FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Seed básico (desenvolvimento)
INSERT INTO areas (id, name, description) VALUES
  ('11111111-0000-0000-0000-000000000001', 'Produção',   'Linha de produção principal'),
  ('11111111-0000-0000-0000-000000000002', 'Manutenção', 'Manutenção industrial');

-- Senha: admin123  (bcrypt $2a$10$...)
INSERT INTO users (id, name, email, password_hash, role) VALUES
  ('22222222-0000-0000-0000-000000000001', 'Admin Sistema', 'admin@planetaazul.com',
   '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin');

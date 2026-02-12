-- Migration for Project Management (Issues, Projects, etc.)

CREATE TABLE IF NOT EXISTS issue (
    id BIGSERIAL PRIMARY KEY,
    repo_id BIGINT NOT NULL, -- references repository(id)
    number BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT,
    state VARCHAR(20) DEFAULT 'open',
    author_id BIGINT NOT NULL, -- references user(id)
    milestone_id BIGINT,
    assignee_id BIGINT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    closed_at TIMESTAMPTZ,
    UNIQUE(repo_id, number)
);

CREATE TABLE IF NOT EXISTS label (
    id BIGSERIAL PRIMARY KEY,
    repo_id BIGINT NOT NULL, -- references repository(id)
    name VARCHAR(50) NOT NULL,
    color VARCHAR(7) DEFAULT '#ffffff',
    description VARCHAR(255),
    UNIQUE(repo_id, name)
);

CREATE TABLE IF NOT EXISTS issue_label (
    issue_id BIGINT NOT NULL REFERENCES issue(id) ON DELETE CASCADE,
    label_id BIGINT NOT NULL REFERENCES label(id) ON DELETE CASCADE,
    PRIMARY KEY(issue_id, label_id)
);

CREATE TABLE IF NOT EXISTS milestone (
    id BIGSERIAL PRIMARY KEY,
    repo_id BIGINT NOT NULL, -- references repository(id)
    title VARCHAR(255) NOT NULL,
    description TEXT,
    state VARCHAR(20) DEFAULT 'open',
    due_on TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(repo_id, title)
);

CREATE TABLE IF NOT EXISTS comment (
    id BIGSERIAL PRIMARY KEY,
    issue_id BIGINT REFERENCES issue(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL, -- references user(id)
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS project (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL, -- references user(id)
    name VARCHAR(255) NOT NULL,
    description TEXT,
    state VARCHAR(20) DEFAULT 'open',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS project_column (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    position INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS project_card (
    id BIGSERIAL PRIMARY KEY,
    column_id BIGINT NOT NULL REFERENCES project_column(id) ON DELETE CASCADE,
    content_url VARCHAR(500),
    note TEXT,
    position INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

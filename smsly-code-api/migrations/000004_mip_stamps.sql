CREATE TABLE mip_stamp (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id         BIGINT NOT NULL,
    commit_sha      VARCHAR(40) NOT NULL,
    merkle_root     VARCHAR(64) NOT NULL,
    tree_hash       VARCHAR(64) NOT NULL,
    author_id       BIGINT NOT NULL,
    parent_stamp_id UUID REFERENCES mip_stamp(id),
    signature       TEXT NOT NULL,
    verified        BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(repo_id, commit_sha)
);

CREATE INDEX idx_mip_stamp_repo ON mip_stamp(repo_id);
CREATE INDEX idx_mip_stamp_commit ON mip_stamp(commit_sha);

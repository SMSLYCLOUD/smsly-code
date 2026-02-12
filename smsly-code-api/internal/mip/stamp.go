package mip

import (
	"time"

	"github.com/google/uuid"
)

// MIPStamp represents a cryptographic integrity stamp for a Git commit.
type MIPStamp struct {
	ID            uuid.UUID  `json:"id"`
	RepoID        int64      `json:"repo_id"`
	CommitSHA     string     `json:"commit_sha"`
	MerkleRoot    string     `json:"merkle_root"`
	TreeHash      string     `json:"tree_hash"`
	AuthorID      int64      `json:"author_id"`
	ParentStampID *uuid.UUID `json:"parent_stamp_id"`
	Signature     string     `json:"signature"`
	Verified      bool       `json:"verified"`
	CreatedAt     time.Time  `json:"created_at"`
}

// CreateStamp initializes a new MIPStamp.
// Note: It does NOT sign the stamp. Use Sign() for that.
func CreateStamp(repoID int64, commitSHA, merkleRoot, treeHash string, authorID int64, parentStampID *uuid.UUID) *MIPStamp {
	return &MIPStamp{
		ID:            uuid.New(),
		RepoID:        repoID,
		CommitSHA:     commitSHA,
		MerkleRoot:    merkleRoot,
		TreeHash:      treeHash,
		AuthorID:      authorID,
		ParentStampID: parentStampID,
		CreatedAt:     time.Now().UTC(),
		Verified:      false,
	}
}

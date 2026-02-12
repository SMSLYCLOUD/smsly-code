package mip

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// mipPayload defines the exact fields used for signature generation/verification.
type mipPayload struct {
	CommitSHA     string     `json:"commit_sha"`
	MerkleRoot    string     `json:"merkle_root"`
	TreeHash      string     `json:"tree_hash"`
	AuthorID      int64      `json:"author_id"`
	ParentStampID *uuid.UUID `json:"parent_stamp_id"`
	Timestamp     time.Time  `json:"timestamp"`
}

// Sign generates an Ed25519 signature for the stamp using the provided private key.
// It populates the Signature field of the stamp.
func (s *MIPStamp) Sign(privateKey ed25519.PrivateKey) error {
	if len(privateKey) != ed25519.PrivateKeySize {
		return errors.New("invalid private key size")
	}

	payload := mipPayload{
		CommitSHA:     s.CommitSHA,
		MerkleRoot:    s.MerkleRoot,
		TreeHash:      s.TreeHash,
		AuthorID:      s.AuthorID,
		ParentStampID: s.ParentStampID,
		Timestamp:     s.CreatedAt, // Ensure we use CreatedAt as "timestamp"
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload for signing: %w", err)
	}

	signature := ed25519.Sign(privateKey, data)
	s.Signature = hex.EncodeToString(signature)
	return nil
}

// VerifySignature verifies the Ed25519 signature against the provided public key.
func (s *MIPStamp) VerifySignature(publicKey ed25519.PublicKey) (bool, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return false, errors.New("invalid public key size")
	}

	if s.Signature == "" {
		return false, errors.New("signature is missing")
	}

	sigBytes, err := hex.DecodeString(s.Signature)
	if err != nil {
		return false, fmt.Errorf("invalid signature hex: %w", err)
	}

	if len(sigBytes) != ed25519.SignatureSize {
		return false, errors.New("invalid signature size")
	}

	payload := mipPayload{
		CommitSHA:     s.CommitSHA,
		MerkleRoot:    s.MerkleRoot,
		TreeHash:      s.TreeHash,
		AuthorID:      s.AuthorID,
		ParentStampID: s.ParentStampID,
		Timestamp:     s.CreatedAt,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal payload for verification: %w", err)
	}

	valid := ed25519.Verify(publicKey, data, sigBytes)
	return valid, nil
}

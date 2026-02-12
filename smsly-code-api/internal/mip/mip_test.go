package mip

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/google/uuid"
)

func generateKeys(t *testing.T) (ed25519.PublicKey, ed25519.PrivateKey) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate keys: %v", err)
	}
	return pub, priv
}

func TestCreateStamp(t *testing.T) {
	stamp := CreateStamp(1, "abc", "root", "tree", 10, nil)
	if stamp.RepoID != 1 {
		t.Errorf("expected RepoID 1, got %d", stamp.RepoID)
	}
	if stamp.CommitSHA != "abc" {
		t.Errorf("expected CommitSHA abc, got %s", stamp.CommitSHA)
	}
	if stamp.Verified {
		t.Errorf("expected Verified false, got true")
	}
	if stamp.CreatedAt.IsZero() {
		t.Errorf("expected CreatedAt to be set")
	}
}

func TestComputeMerkleRoot_Empty(t *testing.T) {
	root, err := ComputeMerkleRoot(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root == "" {
		t.Errorf("expected non-empty root for empty list")
	}
}

func TestComputeMerkleRoot_SingleFile(t *testing.T) {
	files := []FileEntry{{Path: "a", Hash: "h1"}}
	root, err := ComputeMerkleRoot(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root == "" {
		t.Errorf("expected non-empty root")
	}
}

func TestComputeMerkleRoot_Determinism(t *testing.T) {
	files1 := []FileEntry{{Path: "b", Hash: "h2"}, {Path: "a", Hash: "h1"}}
	files2 := []FileEntry{{Path: "a", Hash: "h1"}, {Path: "b", Hash: "h2"}}

	root1, _ := ComputeMerkleRoot(files1)
	root2, _ := ComputeMerkleRoot(files2)

	if root1 != root2 {
		t.Errorf("expected deterministic root, got %s vs %s", root1, root2)
	}
}

func TestComputeMerkleRoot_DifferentContent(t *testing.T) {
	files1 := []FileEntry{{Path: "a", Hash: "h1"}}
	files2 := []FileEntry{{Path: "a", Hash: "h2"}}

	root1, _ := ComputeMerkleRoot(files1)
	root2, _ := ComputeMerkleRoot(files2)

	if root1 == root2 {
		t.Errorf("expected different roots for different content")
	}
}

func TestSignVerify_Valid(t *testing.T) {
	pub, priv := generateKeys(t)
	stamp := CreateStamp(1, "sha", "root", "tree", 1, nil)

	err := stamp.Sign(priv)
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	if stamp.Signature == "" {
		t.Errorf("signature is empty")
	}

	valid, err := VerifyStamp(stamp, pub)
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if !valid {
		t.Errorf("expected signature to be valid")
	}
}

func TestSignVerify_Tampered(t *testing.T) {
	pub, priv := generateKeys(t)
	stamp := CreateStamp(1, "sha", "root", "tree", 1, nil)
	stamp.Sign(priv)

	// Tamper
	stamp.CommitSHA = "tampered"
	valid, _ := VerifyStamp(stamp, pub)
	if valid {
		t.Errorf("expected signature to be invalid after tampering SHA")
	}

	stamp.CommitSHA = "sha" // restore
	stamp.AuthorID = 999
	valid, _ = VerifyStamp(stamp, pub)
	if valid {
		t.Errorf("expected signature to be invalid after tampering AuthorID")
	}
}

func TestSignVerify_InvalidKey(t *testing.T) {
	_, priv := generateKeys(t)
	pub2, _ := generateKeys(t)
	stamp := CreateStamp(1, "sha", "root", "tree", 1, nil)
	stamp.Sign(priv)

	valid, _ := VerifyStamp(stamp, pub2) // Verify with wrong key
	if valid {
		t.Errorf("expected signature to be invalid with wrong key")
	}
}

func TestSignVerify_MissingSignature(t *testing.T) {
	pub, _ := generateKeys(t)
	stamp := CreateStamp(1, "sha", "root", "tree", 1, nil)

	valid, err := VerifyStamp(stamp, pub)
	if valid {
		t.Errorf("expected invalid for missing signature")
	}
	if err == nil {
		t.Errorf("expected error for missing signature")
	}
}

func TestSignVerify_InvalidSignatureFormat(t *testing.T) {
	pub, _ := generateKeys(t)
	stamp := CreateStamp(1, "sha", "root", "tree", 1, nil)
	stamp.Signature = "invalid-hex"

	valid, err := VerifyStamp(stamp, pub)
	if valid {
		t.Errorf("expected invalid")
	}
	if err == nil {
		t.Errorf("expected error for bad hex")
	}
}

func TestVerifyChain_SingleValid(t *testing.T) {
	pub, priv := generateKeys(t)
	s1 := CreateStamp(1, "c1", "r1", "t1", 1, nil)
	s1.Sign(priv)

	res := VerifyChain([]*MIPStamp{s1}, pub)
	if !res.Valid {
		t.Errorf("expected valid chain for single stamp: %s", res.Error)
	}
}

func TestVerifyChain_TwoValid(t *testing.T) {
	pub, priv := generateKeys(t)
	s1 := CreateStamp(1, "c1", "r1", "t1", 1, nil)
	s1.CreatedAt = time.Now().Add(-1 * time.Hour)
	s1.Sign(priv)

	s2 := CreateStamp(1, "c2", "r2", "t2", 1, &s1.ID)
	s2.CreatedAt = time.Now()
	s2.Sign(priv)

	res := VerifyChain([]*MIPStamp{s1, s2}, pub)
	if !res.Valid {
		t.Errorf("expected valid chain: %s", res.Error)
	}
}

func TestVerifyChain_UnsortedInput(t *testing.T) {
	pub, priv := generateKeys(t)
	s1 := CreateStamp(1, "c1", "r1", "t1", 1, nil)
	s1.CreatedAt = time.Now().Add(-1 * time.Hour)
	s1.Sign(priv)

	s2 := CreateStamp(1, "c2", "r2", "t2", 1, &s1.ID)
	s2.CreatedAt = time.Now()
	s2.Sign(priv)

	// Pass in reverse order (s2, s1)
	res := VerifyChain([]*MIPStamp{s2, s1}, pub)
	if !res.Valid {
		t.Errorf("expected valid chain even if unsorted: %s", res.Error)
	}
}

func TestVerifyChain_BrokenLink(t *testing.T) {
	pub, priv := generateKeys(t)
	s1 := CreateStamp(1, "c1", "r1", "t1", 1, nil)
	s1.CreatedAt = time.Now().Add(-1 * time.Hour)
	s1.Sign(priv)

	// s2 points to random parent, not s1
	randomID := uuid.New()
	s2 := CreateStamp(1, "c2", "r2", "t2", 1, &randomID)
	s2.CreatedAt = time.Now()
	s2.Sign(priv)

	res := VerifyChain([]*MIPStamp{s1, s2}, pub)
	if res.Valid {
		t.Errorf("expected broken chain")
	}
	if res.BrokenAt == nil || *res.BrokenAt != s2.ID {
		t.Errorf("expected broken at s2")
	}
}

func TestVerifyChain_MissingParent(t *testing.T) {
	pub, priv := generateKeys(t)
	s1 := CreateStamp(1, "c1", "r1", "t1", 1, nil)
	s1.CreatedAt = time.Now().Add(-1 * time.Hour)
	s1.Sign(priv)

	// s2 has NO parent, but is second in chain
	s2 := CreateStamp(1, "c2", "r2", "t2", 1, nil)
	s2.CreatedAt = time.Now()
	s2.Sign(priv)

	res := VerifyChain([]*MIPStamp{s1, s2}, pub)
	if res.Valid {
		t.Errorf("expected invalid chain due to missing parent on s2")
	}
}

func TestVerifyChain_TamperedSignatureInChain(t *testing.T) {
	pub, priv := generateKeys(t)
	s1 := CreateStamp(1, "c1", "r1", "t1", 1, nil)
	s1.Sign(priv)

	s2 := CreateStamp(1, "c2", "r2", "t2", 1, &s1.ID)
	s2.Sign(priv)

	// Tamper s1 signature
	s1.Signature = hex.EncodeToString([]byte("bad"))

	res := VerifyChain([]*MIPStamp{s1, s2}, pub)
	if res.Valid {
		t.Errorf("expected invalid chain due to bad signature on s1")
	}
	if res.BrokenAt == nil || *res.BrokenAt != s1.ID {
		t.Errorf("expected broken at s1")
	}
}

func TestVerifyChain_Empty(t *testing.T) {
	pub, _ := generateKeys(t)
	res := VerifyChain([]*MIPStamp{}, pub)
	if !res.Valid {
		t.Errorf("expected valid for empty chain")
	}
}

func TestMerkle_OddNodes(t *testing.T) {
	files := []FileEntry{
		{Path: "a", Hash: "h1"},
		{Path: "b", Hash: "h2"},
		{Path: "c", Hash: "h3"},
	}
	root, err := ComputeMerkleRoot(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root == "" {
		t.Errorf("expected root")
	}
}

func TestVerifyChain_ThreeValid(t *testing.T) {
	pub, priv := generateKeys(t)
	s1 := CreateStamp(1, "c1", "r1", "t1", 1, nil)
	s1.CreatedAt = time.Now().Add(-2 * time.Hour)
	s1.Sign(priv)

	s2 := CreateStamp(1, "c2", "r2", "t2", 1, &s1.ID)
	s2.CreatedAt = time.Now().Add(-1 * time.Hour)
	s2.Sign(priv)

	s3 := CreateStamp(1, "c3", "r3", "t3", 1, &s2.ID)
	s3.CreatedAt = time.Now()
	s3.Sign(priv)

	res := VerifyChain([]*MIPStamp{s1, s2, s3}, pub)
	if !res.Valid {
		t.Errorf("expected valid chain: %s", res.Error)
	}
}

func TestVerifyChain_BrokenMiddle(t *testing.T) {
	pub, priv := generateKeys(t)
	s1 := CreateStamp(1, "c1", "r1", "t1", 1, nil)
	s1.CreatedAt = time.Now().Add(-2 * time.Hour)
	s1.Sign(priv)

	// s2 points to s1
	s2 := CreateStamp(1, "c2", "r2", "t2", 1, &s1.ID)
	s2.CreatedAt = time.Now().Add(-1 * time.Hour)
	s2.Sign(priv)

	// s3 points to s1 (skip s2) - BROKEN
	s3 := CreateStamp(1, "c3", "r3", "t3", 1, &s1.ID)
	s3.CreatedAt = time.Now()
	s3.Sign(priv)

	res := VerifyChain([]*MIPStamp{s1, s2, s3}, pub)
	if res.Valid {
		t.Errorf("expected broken chain")
	}
	if res.BrokenAt == nil || *res.BrokenAt != s3.ID {
		t.Errorf("expected broken at s3")
	}
}

func TestVerifyChain_BrokenEnd(t *testing.T) {
	pub, priv := generateKeys(t)
	s1 := CreateStamp(1, "c1", "r1", "t1", 1, nil)
	s1.CreatedAt = time.Now().Add(-2 * time.Hour)
	s1.Sign(priv)

	s2 := CreateStamp(1, "c2", "r2", "t2", 1, &s1.ID)
	s2.CreatedAt = time.Now().Add(-1 * time.Hour)
	s2.Sign(priv)

	// s3 points to nil (orphan) - BROKEN
	s3 := CreateStamp(1, "c3", "r3", "t3", 1, nil)
	s3.CreatedAt = time.Now()
	s3.Sign(priv)

	res := VerifyChain([]*MIPStamp{s1, s2, s3}, pub)
	if res.Valid {
		t.Errorf("expected broken chain")
	}
}

func TestMIPStamp_JSON(t *testing.T) {
	stamp := CreateStamp(1, "sha", "root", "tree", 1, nil)
	if stamp.ID == uuid.Nil {
		t.Error("ID should be generated")
	}
	if stamp.Verified {
		t.Error("Verified should be false by default")
	}
}

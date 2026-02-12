package mip

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

// FileEntry represents a file in the repository for Merkle tree computation.
type FileEntry struct {
	Path string
	Hash string
}

// ComputeMerkleRoot calculates the SHA-256 Merkle root for a list of files.
// It sorts the files by path to ensure determinism.
func ComputeMerkleRoot(files []FileEntry) (string, error) {
	if len(files) == 0 {
		// Return empty hash or zero hash? Let's use empty string SHA-256 for empty tree.
		h := sha256.Sum256([]byte(""))
		return hex.EncodeToString(h[:]), nil
	}

	// 1. Sort files by path to ensure deterministic order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	// 2. Create leaf nodes
	var nodes []string
	for _, f := range files {
		// Leaf hash = SHA256(path + ":" + hash)
		data := fmt.Sprintf("%s:%s", f.Path, f.Hash)
		h := sha256.Sum256([]byte(data))
		nodes = append(nodes, hex.EncodeToString(h[:]))
	}

	// 3. Build tree
	for len(nodes) > 1 {
		var nextLevel []string

		for i := 0; i < len(nodes); i += 2 {
			if i+1 < len(nodes) {
				// Pair: SHA256(left + right)
				combined := nodes[i] + nodes[i+1]
				h := sha256.Sum256([]byte(combined))
				nextLevel = append(nextLevel, hex.EncodeToString(h[:]))
			} else {
				// Orphan: SHA256(node + node) or just carry over?
				// Bitcoin duplicates the last node.
				// Git uses tree objects differently.
				// Let's assume duplication for simplicity/robustness if odd number.
				// Or carry over.
				// "build binary tree" usually implies pairing.
				// Let's duplicate the last one if odd, common in Merkle implementations.
				combined := nodes[i] + nodes[i]
				h := sha256.Sum256([]byte(combined))
				nextLevel = append(nextLevel, hex.EncodeToString(h[:]))
			}
		}
		nodes = nextLevel
	}

	return nodes[0], nil
}

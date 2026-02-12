package mip

import (
	"crypto/ed25519"
	"fmt"
	"sort"
)

// VerifyStamp verifies the signature of a single MIPStamp.
func VerifyStamp(stamp *MIPStamp, pubKey ed25519.PublicKey) (bool, error) {
	if stamp == nil {
		return false, fmt.Errorf("nil stamp")
	}
	return stamp.VerifySignature(pubKey)
}

// VerifyChain verifies a sequence of MIPStamps.
// The stamps are assumed to be a continuous chain, but their order in the slice is not guaranteed.
// This function will sort them by CreatedAt before verification.
// Returns a ChainVerification result.
func VerifyChain(stamps []*MIPStamp, pubKey ed25519.PublicKey) ChainVerification {
	if len(stamps) == 0 {
		return ChainVerification{Valid: true}
	}

	// 1. Sort stamps by CreatedAt to establish expected chain order
	sortedStamps := make([]*MIPStamp, len(stamps))
	copy(sortedStamps, stamps)
	sort.Slice(sortedStamps, func(i, j int) bool {
		return sortedStamps[i].CreatedAt.Before(sortedStamps[j].CreatedAt)
	})

	for i, stamp := range sortedStamps {
		// 2. Verify individual signature
		valid, err := VerifyStamp(stamp, pubKey)
		if !valid || err != nil {
			msg := "invalid signature"
			if err != nil {
				msg = err.Error()
			}
			return ChainVerification{
				Valid:    false,
				BrokenAt: &stamp.ID,
				Error:    fmt.Sprintf("Stamp %s: %s", stamp.ID, msg),
			}
		}

		// 3. Verify parent linkage (skip for the very first stamp in time)
		if i > 0 {
			prev := sortedStamps[i-1]
			if stamp.ParentStampID == nil {
				return ChainVerification{
					Valid:    false,
					BrokenAt: &stamp.ID,
					Error:    fmt.Sprintf("Stamp %s has no parent, but is not the first in chain (prev: %s)", stamp.ID, prev.ID),
				}
			}
			if *stamp.ParentStampID != prev.ID {
				return ChainVerification{
					Valid:    false,
					BrokenAt: &stamp.ID,
					Error:    fmt.Sprintf("Stamp %s parent mismatch: expected %s, got %s", stamp.ID, prev.ID, *stamp.ParentStampID),
				}
			}
		}
	}

	return ChainVerification{Valid: true}
}

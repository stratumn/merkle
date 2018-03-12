// Copyright 2017 Stratumn SAS. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// MerkleNodeHashes contains a left, right, and parent hash.
type MerkleNodeHashes struct {
	Left   []byte `json:"left"`
	Right  []byte `json:"right"`
	Parent []byte `json:"parent"`
}

// Path contains the necessary hashes to go from a leaf to a Merkle root.
type Path []MerkleNodeHashes

// Validate validates the integrity of a hash triplet.
func (h MerkleNodeHashes) Validate() error {
	hash := sha256.New()

	if _, err := hash.Write(h.Left); err != nil {
		return err
	}
	if _, err := hash.Write(h.Right); err != nil {
		return err
	}

	expected := hash.Sum(nil)

	if bytes.Compare(h.Parent, expected) != 0 {
		var (
			got  = hex.EncodeToString(h.Parent)
			want = hex.EncodeToString(expected)
		)
		return fmt.Errorf("unexpected parent hash got %q want %q", got, want)
	}

	return nil
}

// Validate validates the integrity of a Merkle path.
func (p Path) Validate() error {
	for i, h := range p {
		if err := h.Validate(); err != nil {
			return err
		}

		if i < len(p)-1 {
			up := p[i+1]

			if bytes.Compare(h.Parent, up.Left) != 0 && bytes.Compare(h.Parent, up.Right) != 0 {
				var (
					e  = hex.EncodeToString(h.Parent)
					a1 = hex.EncodeToString(up.Left)
					a2 = hex.EncodeToString(up.Right)
				)
				return fmt.Errorf("could not find parent hash %q, got %q and %q", e, a1, a2)
			}
		}
	}

	return nil
}

// JSONMerkleNodeHashes is used to Marshal/Unmarshal MerkleNodeHashes type with
// hex representation.
type JSONMerkleNodeHashes struct {
	Left   string `json:"left"`
	Right  string `json:"right"`
	Parent string `json:"parent"`
}

// MarshalJSON implements encoding/json.Marshaler.MarshalJSON.
func (h *MerkleNodeHashes) MarshalJSON() ([]byte, error) {
	return json.Marshal(JSONMerkleNodeHashes{
		Left:   hex.EncodeToString(h.Left),
		Right:  hex.EncodeToString(h.Right),
		Parent: hex.EncodeToString(h.Parent),
	})
}

// UnmarshalJSON implements encoding/json.Unmarshaler.UnmarshalJSON.
func (h *MerkleNodeHashes) UnmarshalJSON(data []byte) error {
	var j JSONMerkleNodeHashes
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	var err error
	if h.Left, err = hex.DecodeString(j.Left); err != nil {
		return err
	}
	if h.Right, err = hex.DecodeString(j.Right); err != nil {
		return err
	}
	if h.Parent, err = hex.DecodeString(j.Parent); err != nil {
		return err
	}
	return nil
}

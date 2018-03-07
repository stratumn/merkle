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

// Package types defines common types.
package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Bytes32Size is the size of a 32-byte long byte array.
const Bytes32Size = 32

// Bytes32Zero is the default value for a 32-byte long byte array.
var Bytes32Zero = &Bytes32{}

// Bytes32 is a 32-byte long byte array.
type Bytes32 [Bytes32Size]byte

// String returns a hex encoded string.
func (b *Bytes32) String() string {
	return hex.EncodeToString(b[:])
}

// Unstring sets the value from a hex encoded string.
func (b *Bytes32) Unstring(src string) error {
	buf, err := hex.DecodeString(src)
	if err != nil {
		return err
	}
	if n := len(buf); n != Bytes32Size {
		return fmt.Errorf("invalid Bytes32 size got %d want %d", n, Bytes32Size)
	}

	copy(b[:], buf)
	return nil
}

// MarshalJSON implements encoding/json.Marshaler.MarshalJSON.
func (b *Bytes32) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

// UnmarshalJSON implements encoding/json.Unmarshaler.UnmarshalJSON.
func (b *Bytes32) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return b.Unstring(s)
}

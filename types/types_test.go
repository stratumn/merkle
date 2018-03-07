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

package types_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stratumn/merkle/types"
)

func TestBytes32String(t *testing.T) {
	str := "1234567890123456789012345678901234567890123456789012345678901234"
	buf, _ := hex.DecodeString(str)
	var b types.Bytes32
	copy(b[:], buf)

	if got, want := b.String(), str; got != want {
		t.Errorf("b.String() = %q want %q", got, want)
	}
}

func TestBytes32Unstring(t *testing.T) {
	str := "1234567890123456789012345678901234567890123456789012345678901234"
	var b types.Bytes32
	if err := b.Unstring(str); err != nil {
		t.Fatalf("b.Unstring(): err: %s", err)
	}
	if got, want := b.String(), str; got != want {
		t.Errorf("b.String() = %q want %q", got, want)
	}
}

func TestBytes32Unstring_invalidHex(t *testing.T) {
	var b types.Bytes32
	if err := b.Unstring("123y567890123456789012345678901234567890123456789012345678901234"); err == nil {
		t.Error("b.Unstring(): err = nil want Error")
	}
}

func TestBytes32Unstring_invalidSize(t *testing.T) {
	var b types.Bytes32
	if err := b.Unstring("17890123456789012345678901234567890123456789012345"); err == nil {
		t.Error("b.Unstring(): err = nil want Error")
	}
}

func TestBytes32MarshalJSON(t *testing.T) {
	str := "1234567890123456789012345678901234567890123456789012345678901234"
	buf, _ := hex.DecodeString(str)
	var b types.Bytes32
	copy(b[:], buf)
	marshalled, err := json.Marshal(&b)
	if err != nil {
		t.Fatalf("json.Marshal(): err: %s", err)
	}

	if got, want := string(marshalled), fmt.Sprintf(`"%s"`, str); got != want {
		t.Errorf("b.MarshalJSON() = %q want %q", got, want)
	}
}

func TestBytes32UnmarshalJSON(t *testing.T) {
	str := "1234567890123456789012345678901234567890123456789012345678901234"
	marshalled := fmt.Sprintf(`"%s"`, str)
	var b types.Bytes32
	err := json.Unmarshal([]byte(marshalled), &b)
	if err != nil {
		t.Fatalf("json.Unmarshal(): err: %s", err)
	}

	if got, want := b.String(), str; got != want {
		t.Errorf("b.UnmarshalJSON() = %q want %q", got, want)
	}
}

func TestBytes32UnmarshalJSON_invalidStr(t *testing.T) {
	marshalled, err := json.Marshal([]string{"test"})
	if err != nil {
		t.Fatalf("json.Marshal(): err: %s", err)
	}
	var b types.Bytes32
	err = json.Unmarshal([]byte(marshalled), &b)
	if err == nil {
		t.Error("json.Unmarshal(): err = nil want Error")
	}
}

func TestBytes32UnmarshalJSON_invalidHex(t *testing.T) {
	str := "t234567890123456789012345678901234567890123456789012345678901234"
	marshalled := fmt.Sprintf(`"%s"`, str)
	var b types.Bytes32
	err := json.Unmarshal([]byte(marshalled), &b)
	if err == nil {
		t.Error("json.Unmarshal(): err = nil want Error")
	}
}

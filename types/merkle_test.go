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
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stratumn/merkle/testutil"
	"github.com/stratumn/merkle/types"
)

func loadPath(filename string, path *types.Path) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, path); err != nil {
		return err
	}
	return nil
}

func TestMerkleNodeHashesValidate_OK(t *testing.T) {
	var (
		left  = testutil.RandomHash()
		right = testutil.RandomHash()
		h     = types.MerkleNodeHashes{Left: left, Right: right}
		hash  = sha256.New()
	)

	if _, err := hash.Write(left); err != nil {
		t.Fatalf("hash.Write(): err: %s", err)
	}
	if _, err := hash.Write(right); err != nil {
		t.Fatalf("hash.Write(): err: %s", err)
	}

	h.Parent = hash.Sum(nil)

	if err := h.Validate(); err != nil {
		t.Errorf("h.Validate(): err: %s", err)
	}
}

func TestMerkleNodeHashesValidate_Error(t *testing.T) {
	h := types.MerkleNodeHashes{
		Left:   testutil.RandomHash(),
		Right:  testutil.RandomHash(),
		Parent: testutil.RandomHash(),
	}
	if err := h.Validate(); err == nil {
		t.Error("h.Validate(): err = nil want Error")
	}
}

func TestPathValidate_OK(t *testing.T) {
	var (
		pathABCDE0 types.Path
		pathABCDE4 types.Path
	)
	if err := loadPath("testdata/path-abcde-0.json", &pathABCDE0); err != nil {
		t.Fatalf("loadPath(): err: %s", err)
	}
	if err := loadPath("testdata/path-abcde-4.json", &pathABCDE4); err != nil {
		t.Fatalf("loadPath(): err: %s", err)
	}

	if err := pathABCDE0.Validate(); err != nil {
		t.Errorf("pathABCDE0.Validate(): err: %s", err)
	}
	if err := pathABCDE4.Validate(); err != nil {
		t.Errorf("pathABCDE4.Validate(): err: %s", err)
	}
}

func TestPathValidate_Error(t *testing.T) {
	var (
		pathInvalid0 types.Path
		pathInvalid1 types.Path
	)
	if err := loadPath("testdata/path-invalid-0.json", &pathInvalid0); err != nil {
		t.Fatalf("loadPath(): err: %s", err)
	}
	if err := loadPath("testdata/path-invalid-1.json", &pathInvalid1); err != nil {
		t.Fatalf("loadPath(): err: %s", err)
	}

	if err := pathInvalid0.Validate(); err == nil {
		t.Error("pathInvalid0.Validate(): err = nil want Error")
	}
	if err := pathInvalid1.Validate(); err == nil {
		t.Error("pathInvalid1.Validate(): err = nil want Error")
	}
}

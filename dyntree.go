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

package merkle

import (
	"crypto/sha256"
	"hash"
	"sync"

	"github.com/stratumn/merkle/types"
)

// DynTreeNode is a node within a DynTree.
type DynTreeNode struct {
	hash   []byte
	left   *DynTreeNode
	right  *DynTreeNode
	parent *DynTreeNode
	height int
}

// Hash returns the hash of the node.
func (n *DynTreeNode) Hash() []byte {
	return n.hash
}

// Left returns the node to the left, if any.
func (n *DynTreeNode) Left() *DynTreeNode {
	return n.left
}

// Right returns the node to the right, if any.
func (n *DynTreeNode) Right() *DynTreeNode {
	return n.right
}

// Parent returns the parent node, if any.
func (n *DynTreeNode) Parent() *DynTreeNode {
	return n.parent
}

func (n *DynTreeNode) rehash(h hash.Hash, a, b []byte, rehashParent bool) {
	h.Reset()

	// Write never returns an error.
	h.Write(a)
	h.Write(b)
	n.hash = h.Sum(nil)

	if rehashParent && n.parent != nil {
		if n.left != nil {
			n.parent.rehash(h, n.left.hash, n.hash, true)
		} else {
			n.parent.rehash(h, n.hash, n.right.hash, true)
		}
	}
}

// DynTree is designed for Merkle trees that can mutate.
// It supports pausing/resuming the computation of hashes, which is useful
// when adding a large number of leaves to the tree to gain more performance
type DynTree struct {
	nodes  []DynTreeNode
	root   *DynTreeNode
	leaves []*DynTreeNode
	height int
	mutex  sync.RWMutex
	hash   hash.Hash
	paused bool
}

// NewDynTree creates a DynTree.
func NewDynTree(initialCap int) *DynTree {
	return &DynTree{
		nodes:  make([]DynTreeNode, 0, initialCap*2-1),
		leaves: make([]*DynTreeNode, 0, initialCap),
		hash:   sha256.New(),
	}
}

// LeavesLen returns the number of leaves. Implements Tree.LeavesLen.
func (t *DynTree) LeavesLen() int {
	return len(t.leaves)
}

// Root returns the Merkle root. Implements Tree.Root.
func (t *DynTree) Root() []byte {
	return t.root.Hash()
}

// Leaf returns the leaf at the specified index. Implements Tree.Leaf.
func (t *DynTree) Leaf(index int) []byte {
	return t.leaves[index].Hash()
}

// Path returns the path of a leaf to the Merkle root. Implements Tree.Path.
func (t *DynTree) Path(index int) types.Path {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if len(t.leaves) < 2 {
		return types.Path{}
	}

	var (
		path  = make(types.Path, t.height)
		node  = t.leaves[index]
		level = 0
	)

	for node.parent != nil {
		if node.left != nil {
			path[level] = types.MerkleNodeHashes{
				Left:   node.left.hash,
				Right:  node.hash,
				Parent: node.parent.hash,
			}
		} else {
			path[level] = types.MerkleNodeHashes{
				Left:   node.hash,
				Right:  node.right.hash,
				Parent: node.parent.hash,
			}
		}

		node = node.parent
		level++
	}

	return path[:level]
}

// Add adds a leaf to the tree.
func (t *DynTree) Add(leaf []byte) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.nodes = append(t.nodes, DynTreeNode{hash: leaf})
	node := &t.nodes[len(t.nodes)-1]
	t.leaves = append(t.leaves, node)

	if t.root == nil {
		t.root = node
	} else {
		left := t.leaves[len(t.leaves)-2]

		for left.parent != nil && left.parent.height == left.height+1 {
			left = left.parent
		}

		t.nodes = append(t.nodes, DynTreeNode{
			left:   left.left,
			parent: left.parent,
			height: left.height + 1,
		})

		parent := &t.nodes[len(t.nodes)-1]
		node.parent, node.left = parent, left
		left.parent, left.right = parent, node

		if left.left != nil {
			left.left.right, left.left = parent, nil
		}

		if parent.parent == nil {
			t.root = parent
		}

		if parent.height > t.height {
			t.height = parent.height
		}

		if !t.paused {
			parent.rehash(t.hash, left.hash, leaf, true)
		}
	}
}

// Update updates a leaf of the tree.
func (t *DynTree) Update(index int, hash []byte) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	node := t.leaves[index]
	node.hash = hash

	if !t.paused {
		if node.left != nil {
			node.parent.rehash(t.hash, node.left.hash, hash, true)
		} else if node.right != nil {
			node.parent.rehash(t.hash, hash, node.right.hash, true)
		}
	}
}

// Pause pauses the computation of hashes.
func (t *DynTree) Pause() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.paused = true
}

// Resume resumes the computation of hashes.
func (t *DynTree) Resume() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.recompute()
	t.paused = false
}

func (t *DynTree) recompute() {
	rows := t.leaves

	for {
		if len(rows) < 1 {
			break
		}

		top := make([]*DynTreeNode, 0, len(rows)/2)
		height := rows[0].height

		for i := 0; i < len(rows); i += 2 {
			node := rows[i]
			if node.parent != nil && node.parent.height == height+1 {
				node.parent.rehash(t.hash, node.hash, node.right.hash, false)
				top = append(top, node.parent)
			}
		}

		rows = top
	}
}

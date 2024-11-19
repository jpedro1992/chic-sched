package topology

import (
	"bytes"
	"fmt"
	"unsafe"
)

// LTree : a logical tree topology
type LTree struct {
	// extends Tree
	Tree
}

// NewLTree : create a new logical tree
//   - returns nil if bad parameters
func NewLTree(tree *Tree) *LTree {
	if tree == nil {
		return nil
	}
	return &LTree{
		Tree: *tree,
	}
}

// PercolateClaimed : set claimed from the leaves up to the root
func (lTree *LTree) PercolateClaimed() {
	lTree.ResetClaimed(false)
	leaves := lTree.GetLeaves()
	for _, leaf := range leaves {
		claimed := 0
		path := leaf.GetPathToRoot()
		for i, node := range path {
			lNode := (*LNode)(unsafe.Pointer(node))
			if i == 0 {
				claimed = lNode.GetClaimed()
			} else {
				lNode.IncClaimed(claimed)
			}
		}
	}
}

// ResetClaimed : reset the claimed value in all nodes in the tree
func (lTree *LTree) ResetClaimed(includeLeaves bool) {
	if lTree.root != nil {
		lRoot := (*LNode)(unsafe.Pointer(lTree.root))
		lRoot.ResetClaimed(includeLeaves)
	}
}

// SetPhysicalClaimed : set claimed values in corresponding pTree as claimed values in lTree
func (lTree *LTree) SetPhysicalClaimed() {
	for _, node := range lTree.GetNodeListBFS() {
		lNode := (*LNode)(unsafe.Pointer(node))
		pNode := lNode.pNode
		if pNode != nil {
			pNode.SetNumClaimed(lNode.GetClaimed())
		}
	}
}

// String : a print out of the logical tree
func (lTree *LTree) String() string {
	var b bytes.Buffer
	b.WriteString("lTree:\n")
	fmt.Fprintf(&b, "%s", &lTree.Tree)
	b.WriteString("\n")

	b.WriteString("lNodes:\n")
	nodes := lTree.GetNodeListBFS()
	for _, node := range nodes {
		lNode := (*LNode)(unsafe.Pointer(node))
		fmt.Fprintf(&b, "%s\n", lNode)
	}
	b.WriteString("\n")
	return b.String()
}

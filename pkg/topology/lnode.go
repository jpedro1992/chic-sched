package topology

import (
	"fmt"
	"unsafe"

	"github.com/ibm/chic-sched/pkg/system"
)

// LNode : a node in a logical tree topology
//   - returns nil if bad parameters
type LNode struct {
	// extends Node
	Node
	// corresponding node in the physical tree
	pNode *PNode
	// count of LEs in the subtree rooted at this node
	count int

	// number of LEs claimed in the subtree rooted at this node
	claimed int
}

// NewLNode : create a new logical node
func NewLNode(pNode *PNode, count int) *LNode {
	if pNode == nil || count < 0 {
		return nil
	}
	return &LNode{
		Node:    *NewNode(&system.Entity{ID: pNode.GetID()}),
		pNode:   pNode,
		count:   count,
		claimed: 0,
	}
}

// GetCount : get the count of the node
func (lNode *LNode) GetCount() int {
	return lNode.count
}

// SetCount : set the count of the node
func (lNode *LNode) SetCount(count int) {
	lNode.count = count
}

// GetClaimed : get the number claimed on the node
func (lNode *LNode) GetClaimed() int {
	return lNode.claimed
}

// SetClaimed : set the number claimed on the node
func (lNode *LNode) SetClaimed(claimed int) {
	lNode.claimed = claimed
}

// IncClaimed : increment the number claimed on the node
func (lNode *LNode) IncClaimed(deltaClaimed int) {
	lNode.claimed += deltaClaimed
}

// ResetClaimed : set claimed to zero in subtree
func (lNode *LNode) ResetClaimed(includeLeaves bool) {
	if !includeLeaves && lNode.IsLeaf() {
		return
	}
	lNode.claimed = 0
	for _, node := range lNode.children {
		lChild := (*LNode)(unsafe.Pointer(node))
		lChild.ResetClaimed(includeLeaves)
	}
}

// String : a print out of the logical node
func (lNode *LNode) String() string {
	return fmt.Sprintf("lNode: ID=%s; count=%d; claimed=%d", lNode.GetID(), lNode.count, lNode.claimed)
}

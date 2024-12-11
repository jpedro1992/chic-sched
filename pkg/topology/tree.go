package topology

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/jpedro1992/chic-sched/pkg/system"
)

// Tree : a basic tree
type Tree struct {
	// the root of the tree (could be null)
	root *Node
}

// NewTree : create a tree
func NewTree(root *Node) *Tree {
	return &Tree{
		root: root,
	}
}

// GetHeight : the height of the tree
func (t *Tree) GetHeight() int {
	if t.root == nil {
		return 0
	}
	return t.root.GetHeight()
}

// GetRoot : the root of the tree
func (t *Tree) GetRoot() *Node {
	return t.root
}

// GetLeaves : the leaves of the tree
func (t *Tree) GetLeaves() []*Node {
	if t.root == nil {
		return make([]*Node, 0)
	}
	return t.root.GetLeaves()
}

// GetLeafIDs : the IDs of the leaves of the tree
func (t *Tree) GetLeafIDs() []string {
	leafIDs := make([]string, 0)
	if t.root != nil {
		for _, n := range t.root.GetLeaves() {
			leafIDs = append(leafIDs, n.GetID())
		}
	}
	return leafIDs
}

// GetLeavesMap : a map of the leaves in the tree, keyed by ID
func (t *Tree) GetLeavesMap() map[string]*Node {
	leavesMap := make(map[string]*Node)
	leaves := t.GetLeaves()
	for _, leaf := range leaves {
		leavesMap[leaf.GetID()] = leaf
	}
	return leavesMap
}

// GetNode : find node in tree with a given ID; null if not found
func (t *Tree) GetNode(nodeID string) *Node {
	// TODO: more efficient search
	allNodes := t.GetNodeListBFS()
	for _, n := range allNodes {
		if n.GetID() == nodeID {
			return n
		}
	}
	return nil
}

// GetNodeListBFS : list of nodes in BFS order
func (t *Tree) GetNodeListBFS() []*Node {
	allNodes := make([]*Node, 0)
	nodeList := make([]*Node, 0)
	if t.root != nil {
		nodeList = append(nodeList, t.root)
	}
	for len(nodeList) > 0 {
		n := nodeList[0]
		nodeList = nodeList[1:]
		allNodes = append(allNodes, n)
		nodeList = append(nodeList, n.GetChildren()...)
	}
	return allNodes
}

// GetNodeIDs : get (sorted) IDs of all nodes in the tree
func (t *Tree) GetNodeIDs() []string {
	nodes := t.GetNodeListBFS()
	nodeIDs := make([]string, len(nodes))
	for i, node := range nodes {
		nodeIDs[i] = node.GetID()
	}
	sort.Strings(nodeIDs)
	return nodeIDs
}

// CopyByLeafIDs : create a copy of tree, only with specified subset of leaves, may return nil
func (t *Tree) CopyByLeafIDs(leafIDs []string) *Tree {
	if t.root == nil {
		return nil
	}

	// initialize
	leavesMap := t.GetLeavesMap()
	var rootCopy *Node
	allNodes := make(map[string]*Node)
	var prevNode *Node

	// go over leaf IDs
	for _, id := range leafIDs {
		leafNode := leavesMap[id]
		// skip non-existing leaf IDs
		if leafNode == nil {
			continue
		}
		// visit (and copy non-copied) nodes from root to leaf node
		path := leafNode.GetPathToRoot()
		for i := len(path) - 1; i >= 0; i-- {
			curNode := path[i]
			id := curNode.GetID()
			if curNodeCopy, exists := allNodes[id]; exists {
				prevNode = curNodeCopy
				continue
			}
			// create Node copy
			curNodeCopy := NewNode(&system.Entity{ID: id})
			curNodeCopy.SetValue(curNode.GetValue())
			allNodes[id] = curNodeCopy
			// check if we are at the root node, otherwise link to parent
			if prevNode == nil {
				rootCopy = curNodeCopy
			} else {
				prevNode.AddChild(curNodeCopy)
			}
			prevNode = curNodeCopy
		}
	}
	return NewTree(rootCopy)
}

type TreeMap map[string]TreeMap

// ToMap : create a tree map
func (t *Tree) ToMap() TreeMap {
	if t.root == nil {
		return nil
	}
	m := make(map[string]TreeMap)
	m[t.root.GetID()] = t.root.ToMap()
	return m
}

// GetLeavesDistanceFrom : get an ordered list by distance (in levels) between all leaf nodes to a given source leaf node
func (t *Tree) GetLeavesDistanceFrom(leafSourceName string) []*NodeValue {
	nodeValue := make([]*NodeValue, 0)
	leafSourceNode := t.GetNode(leafSourceName)
	// return if source is not a leaf node
	if leafSourceNode == nil || !leafSourceNode.IsLeaf() {
		return nodeValue
	}
	// add the source node
	nodeValue = append(nodeValue, &NodeValue{Name: leafSourceName, Value: 0})
	// visit nodes along path to root
	path := leafSourceNode.GetPathToRoot()
	prevNode := leafSourceNode
	for i := 1; i < len(path); i++ {
		for _, child := range path[i].GetChildren() {
			// exclude visited nodes
			if child == prevNode {
				continue
			}
			// get depth of all leaves from child node
			nvc := child.GetLeavesDepth()
			for _, nv := range nvc {
				// augment depth by distance to source leaf node
				nv.Value += i + 1
			}
			nodeValue = append(nodeValue, nvc...)
		}
		prevNode = path[i]
	}
	// order the nodes by distance
	sort.Slice(nodeValue, func(i, j int) bool {
		return nodeValue[i].Value < nodeValue[j].Value
	})
	return nodeValue
}

// String : a print out of the tree
func (t *Tree) String() string {
	var b bytes.Buffer
	if t.root != nil {
		fmt.Fprintf(&b, "%s", t.root)
	} else {
		b.WriteString("null")
	}
	return b.String()
}

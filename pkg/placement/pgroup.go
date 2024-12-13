package placement

import (
	"bytes"
	"fmt"
	"strconv"
	"unsafe"

	"github.com/ibm/chic-sched/pkg/system"
	"github.com/ibm/chic-sched/pkg/topology"
	"github.com/ibm/chic-sched/pkg/util"
)

// PGroup : a placement group
type PGroup struct {
	// extends Entity
	system.Entity
	// size of the group (number of members)
	size int
	// resource demand of a member
	demand *util.Allocation

	// level constraints mapped to IDs
	lcIDs map[string]*LevelConstraint
	// level constraints mapped to levels
	lcs map[int]*LevelConstraint

	// logical tree for placement
	lTree *topology.LTree
	// group of LEs
	leGroup *system.LEGroup
}

// NewPGroup : create a new placement group
//   - returns nil if bad parameters
func NewPGroup(id string, size int, demand *util.Allocation) *PGroup {
	if len(id) == 0 || size < 0 || demand.GetSize() == 0 {
		return nil
	}
	leGroup := system.NewLEGroup(id, size, demand)
	for i := 0; i < size; i++ {
		namei := id + "-vm" + strconv.FormatInt(int64(i), 10)
		lei := system.NewLE(namei, demand)
		leGroup.AddLE(lei)
	}
	return &PGroup{
		Entity:  system.Entity{ID: id},
		size:    size,
		demand:  demand,
		lcIDs:   make(map[string]*LevelConstraint),
		lcs:     make(map[int]*LevelConstraint),
		lTree:   nil,
		leGroup: leGroup,
	}
}

// GetSize : get group size
func (pg *PGroup) GetSize() int {
	return pg.size
}

// GetDemand : get resource demand
func (pg *PGroup) GetDemand() *util.Allocation {
	return pg.demand
}

// AddLevelConstraint : add a level constraint to this PGroup
func (pg *PGroup) AddLevelConstraint(lc *LevelConstraint) {
	if lc != nil {
		id := lc.GetID()
		level := lc.GetLevel()
		// remove old if duplicate ID
		pg.RemoveLevelConstraint(id)
		// remove old if duplicate level
		pg.RemoveLevelConstraintByLevel(level)
		// add new constraint
		pg.lcIDs[id] = lc
		pg.lcs[level] = lc
	}
}

// RemoveLevelConstraint : remove a level constraint from this PGroup (by ID)
func (pg *PGroup) RemoveLevelConstraint(lcID string) {
	if lc := pg.lcIDs[lcID]; lc != nil {
		delete(pg.lcIDs, lcID)
		delete(pg.lcs, lc.GetLevel())
	}
}

// RemoveLevelConstraintByLevel : remove a level constraint from this PGroup (by level number)
func (pg *PGroup) RemoveLevelConstraintByLevel(level int) {
	if lc := pg.lcs[level]; lc != nil {
		delete(pg.lcIDs, lc.GetID())
		delete(pg.lcs, level)
	}
}

// GetLevelConstraint : get the level constraint at a given level
// (default if not specified)
func (pg *PGroup) GetLevelConstraint(level int) *LevelConstraint {
	lc := pg.lcs[level]
	if lc == nil {
		lcd := *DefaultLevelConstraint
		lcd.level = level
		lc = &lcd
	}
	return lc
}

// GetLevelConstraintIDs : get the IDs of the level constraints in this PGroup
func (pg *PGroup) GetLevelConstraintIDs() []string {
	ids := make([]string, len(pg.lcIDs))
	i := 0
	for id := range pg.lcIDs {
		ids[i] = id
		i++
	}
	return ids
}

// GetLTree : get the logical tree for this placement group
//   - returns nil if unplaced
func (pg *PGroup) GetLTree() *topology.LTree {
	return pg.lTree
}

// SetLTree : get the logical tree for this placement group
func (pg *PGroup) SetLTree(lTree *topology.LTree) {
	pg.lTree = lTree
}

// GetLEGroup : get the group of LEs in this placement group
func (pg *PGroup) GetLEGroup() *system.LEGroup {
	return pg.leGroup
}

// IsFullyPlaced : are all members of the group placed
func (pg *PGroup) IsFullyPlaced() bool {
	if pg.lTree != nil {
		lRoot := (*topology.LNode)(unsafe.Pointer(pg.lTree.GetRoot()))
		if lRoot != nil {
			return lRoot.GetCount() == pg.size
		}
	}
	return false
}

// ClaimAll : claim all members of this placement group and allocate them
func (pg *PGroup) ClaimAll(pTree *topology.PTree) bool {
	return pg.Claim(pg.size, pTree)
}

// Claim : claim n members of this placement group and allocate them
func (pg *PGroup) Claim(n int, pTree *topology.PTree) bool {
	lTree := pg.lTree
	leGroup := pg.GetLEGroup()
	if pTree == nil || lTree == nil || leGroup == nil {
		return false
	}
	lLeaves := lTree.GetLeaves()
	if len(lLeaves) == 0 {
		return false
	}

	// get a map of PEs
	serverMap := pTree.GetPEs()

	// allocate LEs
	les := leGroup.GetLEs()
	index := 0
loop:
	for _, node := range lLeaves {
		lNode := (*topology.LNode)(unsafe.Pointer(node))
		lNode.SetClaimed(0)
		peID := lNode.GetID()
		pe := serverMap[peID]
		if pe != nil {
			for i := 0; i < lNode.GetCount(); i++ {
				lei := les[index]
				index++
				if pe.PlaceLE(lei) {
					lNode.IncClaimed(1)
				}
				if index == n {
					break loop
				}
			}
		}
	}
	lTree.PercolateClaimed()
	lTree.SetPhysicalClaimed()
	pTree.PercolateResources()
	return true
}

// UnClaimAll : unclaim all members of this placement group and deallocate them
func (pg *PGroup) UnClaimAll(pTree *topology.PTree) bool {
	lTree := pg.lTree
	leGroup := pg.GetLEGroup()
	if pTree == nil || lTree == nil || leGroup == nil {
		return false
	}

	// deallocate LEs
	les := leGroup.GetLEs()
	for _, lei := range les {
		pei := lei.GetHost()
		if pei != nil {
			pei.UnPlaceLE(lei)
		}
	}
	lTree.ResetClaimed(true)
	pTree.ResetNumClaimed()
	pTree.PercolateResources()
	return true
}

// String : a print out of the placement group
func (pg *PGroup) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "PG: ID=%s; size=%d; demand=%v; lcs=%v", pg.GetID(), pg.size, pg.demand, pg.GetLevelConstraintIDs())
	b.WriteString("\n")
	if pg.lTree == nil {
		b.WriteString("LTree: nil")
	} else {
		fmt.Fprintf(&b, "%v", pg.lTree)
	}
	return b.String()
}

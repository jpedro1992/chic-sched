package system

import (
	"github.com/jpedro1992/chic-sched/pkg/util"
)

// LEGroup : a group of homogeneous LEs
type LEGroup struct {
	// extends Entity
	Entity
	// group size
	size int
	// homogeneous resource demand
	demand *util.Allocation
	// group members
	group map[string]*LE
}

// NewLEGroup : create a new (empty) group of homogeneous LEs
//   - makes a copy of resource demand allocation
//   - returns nil if bad parameters
func NewLEGroup(id string, size int, demand *util.Allocation) *LEGroup {
	if len(id) == 0 || demand == nil || demand.GetSize() == 0 {
		return nil
	}
	return &LEGroup{
		Entity: Entity{ID: id},
		size:   0,
		demand: demand.Clone(),
		group:  make(map[string]*LE),
	}
}

// AddLE : add an LE to this group
func (leg *LEGroup) AddLE(le *LE) bool {
	if le == nil {
		return false
	}
	if _, exists := leg.group[le.GetID()]; exists {
		return false
	}
	leg.group[le.GetID()] = le
	leg.size++
	return true
}

// RemoveLE : remove an LE from this group
func (leg *LEGroup) RemoveLE(leID string) bool {
	if _, exists := leg.group[leID]; !exists {
		return false
	}
	delete(leg.group, leID)
	leg.size--
	return true
}

// GetSize : the current size of the group
func (leg *LEGroup) GetSize() int {
	return leg.size
}

// GetDemand : the resource demand of the group
//   - returns a copy of resource demand allocation
func (leg *LEGroup) GetDemand() *util.Allocation {
	return leg.demand.Clone()
}

// GetLEs : the LE members of the group
//   - returns a copy of resource demand allocation
func (leg *LEGroup) GetLEs() []*LE {
	les := make([]*LE, len(leg.group))
	i := 0
	for _, le := range leg.group {
		les[i] = le
		i++
	}
	return les
}

// UNDONE:

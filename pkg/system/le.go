package system

import (
	"fmt"

	"github.com/ibm/chic-sched/pkg/util"
)

// LE : Logical Entity
type LE struct {
	// extends Entity
	Entity
	// resource demand
	demand *util.Allocation
	// hosting LE (nil if not hosted)
	host *PE
}

// NewLE : create a new LE
//   - makes a copy of resource demand allocation
//   - returns nil if bad parameters
func NewLE(id string, demand *util.Allocation) *LE {
	if len(id) == 0 || demand == nil || demand.GetSize() == 0 {
		return nil
	}
	return &LE{
		Entity: Entity{ID: id},
		demand: demand.Clone(),
		host:   nil,
	}
}

// GetHost : get the host PE of this LE
//   - nil if not hosted
func (le *LE) GetHost() *PE {
	return le.host
}

// SetHost : set the host PE of this LE
func (le *LE) SetHost(pe *PE) {
	le.host = pe
}

// UNDONE:

// String : a print out of the LE
func (le *LE) String() string {
	hostID := "none"
	if le.host != nil {
		hostID = le.host.GetID()
	}
	return fmt.Sprintf("LE: ID=%s; demand=%v; host=%s", le.GetID(), le.demand, hostID)
}

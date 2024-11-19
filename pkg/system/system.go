package system

// System : data center
type System struct {
	// extends Entity
	Entity
	// list of servers
	servers map[string]*PE
}

// NewSystem : create a new System
//   - returns nil if bad parameters
func NewSystem(id string) *System {
	if len(id) == 0 {
		return nil
	}
	system := &System{
		Entity:  Entity{ID: id},
		servers: make(map[string]*PE),
	}
	return system
}

// UNDONE:

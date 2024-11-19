package system

// Entity : basic element unit
type Entity struct {
	// unique id
	ID string

	// TODO: add common attributes

}

// GetID : the unique ID
func (e *Entity) GetID() string {
	return e.ID
}

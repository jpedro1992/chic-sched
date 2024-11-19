package util

// Affinity : the affinity type of a level constraint
type Affinity int

const (
	// spread placement
	Spread = iota
	// pack placement
	Pack
)

// AffinityToString : get the string reprfesentation of affinity
func AffinityToString(a Affinity) string {
	switch a {
	case Spread:
		return "Spread"
	case Pack:
		return "Pack"
	}
	return "Unknown"
}

var (
	// DefaultTreeKind : the default kind attribute of the tree
	DefaultTreeKind string = "TopologyTree"
)

// TopologyTree : JSON topology tree
type TopologyTree struct {
	Kind     string    `json:"kind"`
	MetaData JMetaData `json:"metadata"`
	Spec     JTreeSpec `json:"spec"`
}

// JMetaData : common metada
type JMetaData struct {
	Name string `json:"name"`
}

// JTreeSpec : spec of topology tree
type JTreeSpec struct {
	ResourceNames []string `json:"resource-names"`
	LevelNames    []string `json:"level-names"`
	Tree          TreeSpec `json:"tree"`
}

// TreeSpec : spec for (sub) tree
type TreeSpec struct {
	Level map[string]TreeSpec `json:"level"`
}

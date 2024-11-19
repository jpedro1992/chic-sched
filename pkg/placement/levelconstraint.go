package placement

import (
	"fmt"

	"github.com/ibm/chic-sched/pkg/system"
	"github.com/ibm/chic-sched/pkg/util"
)

// LevelConstraint : a level constraint for placing a placement group
// by dividing the group into partitions at a given level (non-root)
//   - the partition size may be restricted to a range [min, max]
//   - the number of partitions may be set to a given value
//   - affinity: Pack or Spread
//   - isHard: If true (Hard) then Pack and Spread correspond to ranges [n,n] and [1,1], respectively,
//   - where n is the total number to be placed at the given level,
//   - otherwise (Soft), the range is assumed to be [1,n] with the objective of having the partition
//   - size close to n and 1, respectively
//   - range[min,max]: if Soft, then a range for the partition size may be specified
//   - numPartitions: if Soft, then a number of partitions may be specified
//   - (makes more sense to have a number of partitions with Spread affinity)
//   - Both range[min,max] and numPartitions may be specified; they are considered as strict constraints
//   - factor: number allocated at level is multiple of the value of factor
//   - (note that some combinations of parameter values lead to infeasible solutions)
//
// TODO: Add feasibility checking of range and number of partitions across levels
type LevelConstraint struct {
	// extends Entity
	system.Entity
	// level at which constraint applies
	level int
	// type of affinity
	affinity util.Affinity
	// hard or soft
	isHard bool

	// following are optional if soft

	// min value in range
	minRange int
	// max value in range
	maxRange int
	// number of partitions
	numPartitions int
	// factor
	factor int
}

var (
	// DefaultLevelConstraint : used if level constraint is not specified at a given level
	DefaultLevelConstraint *LevelConstraint = NewLevelConstraint("lc-def", 0, util.Pack, false)
)

// NewLevelConstraint : create a new level constraint
//   - returns nil if bad parameters
func NewLevelConstraint(id string, level int, affinity util.Affinity, isHard bool) *LevelConstraint {
	if len(id) == 0 || level < 0 {
		return nil
	}
	return &LevelConstraint{
		Entity:   system.Entity{ID: id},
		level:    level,
		affinity: affinity,
		isHard:   isHard,
		factor:   1,
	}
}

// GetID : the unique ID
func (lc *LevelConstraint) GetID() string {
	return lc.Entity.ID
}

// GetLevel : get the level of the level constraint
func (lc *LevelConstraint) GetLevel() int {
	return lc.level
}

// Affinity : get the affinity of the level constraint
func (lc *LevelConstraint) Affinity() util.Affinity {
	return lc.affinity
}

// IsHard : level constraint is hard or soft
func (lc *LevelConstraint) IsHard() bool {
	return lc.isHard
}

// SetRange : set a [min, max] range
func (lc *LevelConstraint) SetRange(min int, max int) bool {
	if lc.isHard || min <= 0 || max < min {
		return false
	}
	lc.minRange = min
	lc.maxRange = max
	return true
}

// GetRange : get the [min, max] range and a boolean if set correctly
func (lc *LevelConstraint) GetRange() (min int, max int, valid bool) {
	if !lc.isHard && lc.minRange > 0 && lc.maxRange >= lc.minRange {
		return lc.minRange, lc.maxRange, true
	}
	return 0, 0, false
}

// SetNumPartitions : set the number of partitions
func (lc *LevelConstraint) SetNumPartitions(num int) bool {
	if lc.isHard || num <= 0 {
		return false
	}
	lc.numPartitions = num
	return true
}

// GetNumPartitions : get the number of partitions and a boolean if set correctly
func (lc *LevelConstraint) GetNumPartitions() (num int, valid bool) {
	if !lc.isHard && lc.numPartitions > 0 {
		return lc.numPartitions, true
	}
	return 0, false
}

// SetFactor : set a factor (false if < 1)
func (lc *LevelConstraint) SetFactor(factor int) bool {
	if factor < 1 {
		return false
	}
	lc.factor = factor
	return true
}

// GetFactor : get the factor (false if <= 1)
func (lc *LevelConstraint) GetFactor() (int, bool) {
	if lc.factor <= 1 {
		return 1, false
	}
	return lc.factor, true
}

// String : a print out of the level constraint
func (lc *LevelConstraint) String() string {
	s := fmt.Sprintf("LC: ID=%s; level=%d; affinity=%s; isHard=%v; ",
		lc.GetID(), lc.level, util.AffinityToString(lc.affinity), lc.isHard)
	if min, max, ok := lc.GetRange(); ok {
		s += fmt.Sprintf("range=[%d,%d]; ", min, max)
	}
	if num, ok := lc.GetNumPartitions(); ok {
		s += fmt.Sprintf("numPartitions=%d; ", num)
	}
	if factor, ok := lc.GetFactor(); ok {
		s += fmt.Sprintf("factor=%d; ", factor)
	}
	return s
}

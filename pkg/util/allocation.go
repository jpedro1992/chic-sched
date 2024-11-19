package util

import (
	"bytes"
	"fmt"
	"math"
)

// Allocation : an allocation of an (ordered) array of resources
// (names of resources are left out for efficiency)
type Allocation struct {
	// values of the allocation
	x []int
}

// NewAllocation : create an empty allocation of a given size (length)
func NewAllocation(size int) (*Allocation, error) {
	if size < 0 {
		return nil, fmt.Errorf("invalid size %d", size)
	}
	return &Allocation{
		x: make([]int, size),
	}, nil
}

// NewAllocationCopy : create an allocation given an array of values
func NewAllocationCopy(value []int) (*Allocation, error) {
	a, err := NewAllocation(len(value))
	if err != nil {
		return nil, err
	}
	copy(a.x, value)
	return a, nil
}

// GetSize : get the size (length) of the values array
func (a *Allocation) GetSize() int {
	return len(a.x)
}

// GetValue : get the array of values
func (a *Allocation) GetValue() []int {
	return a.x
}

// SetValue : set the array of values (overwites previous values)
func (a *Allocation) SetValue(value []int) {
	a.x = make([]int, len(value))
	copy(a.x, value)
}

// Clone : create a copy
func (a *Allocation) Clone() *Allocation {
	alloc, _ := NewAllocationCopy(a.x)
	return alloc
}

// Add : add another allocation to this one (false if unequal lengths)
func (a *Allocation) Add(other *Allocation) bool {
	if !a.SameSize(other) {
		return false
	}
	v := other.GetValue()
	for i := 0; i < len(a.x); i++ {
		a.x[i] += v[i]
	}
	return true
}

// Subtract : subtract another allocation to this one (false if unequal lengths)
func (a *Allocation) Subtract(other *Allocation) bool {
	if !a.SameSize(other) {
		return false
	}
	v := other.GetValue()
	for i := 0; i < len(a.x); i++ {
		a.x[i] -= v[i]
	}
	return true
}

// Divide : divide this allocation by another allocation (false if unequal lengths)
func (a *Allocation) Divide(other *Allocation) (*Allocation, error) {
	if !a.SameSize(other) {
		return nil, fmt.Errorf("allocations have different sizes")
	}
	r := make([]int, a.GetSize())
	v := other.GetValue()
	for i := 0; i < len(a.x); i++ {
		if v[i] == 0 {
			r[i] = math.MaxInt32
		} else {
			r[i] = a.x[i] / v[i]
		}
	}
	return NewAllocationCopy(r)
}

// Scale : multiply elements by a given value
func (a *Allocation) Scale(value int) {
	for i := 0; i < len(a.x); i++ {
		a.x[i] *= value
	}
}

// Minimum : min value in the allocation array
func (a *Allocation) Minimum() int {
	min := 0
	if n := len(a.x); n > 0 {
		min = a.x[0]
		for i := 1; i < n; i++ {
			min = Min(min, a.x[i])
		}
	}
	return min
}

// Maximum : max value in the allocation array
func (a *Allocation) Maximum() int {
	max := 0
	if n := len(a.x); n > 0 {
		max = a.x[0]
		for i := 1; i < n; i++ {
			max = Max(max, a.x[i])
		}
	}
	return max
}

// Fit : check if this allocation fits on an entity with a given capacity
// and already allocated values (false if unequal lengths)
func (a *Allocation) Fit(allocated *Allocation, capacity *Allocation) bool {
	available := capacity.Clone()
	if available.Subtract(allocated) {
		return a.LessOrEqual(available)
	}
	return false
}

// NumberToFit : number of this allocation fitting on an entity with a given capacity
// and already allocated values (0 if unequal lengths)
func (a *Allocation) NumberToFit(allocated *Allocation, capacity *Allocation) int {
	available := capacity.Clone()
	if available.Subtract(allocated) {
		if result, err := available.Divide(a); err == nil {
			return result.Minimum()
		}
	}
	return 0
}

// SameSize : check if same size (length) as another allocation
func (a *Allocation) SameSize(other *Allocation) bool {
	return a.GetSize() == other.GetSize()
}

// IsZero : check if values are zeros (all element values)
func (a *Allocation) IsZero() bool {
	for i := 0; i < len(a.x); i++ {
		if a.x[i] != 0 {
			return false
		}
	}
	return true
}

// SetZero : set of all element values to zeros
func (a *Allocation) SetZero() {
	for i := 0; i < len(a.x); i++ {
		a.x[i] = 0
	}
}

// Equal : check if equals another allocation (false if unequal lengths)
func (a *Allocation) Equal(other *Allocation) bool {
	return a.comp(other, false)
}

// LessOrEqual : check if less or equal to another allocation (false if unequal lengths)
func (a *Allocation) LessOrEqual(other *Allocation) bool {
	return a.comp(other, true)
}

// comp : compare this to another allocation (false if unequal lengths);
// lessOrEqual: true => 'less or equal'; false => 'equal'
func (a *Allocation) comp(other *Allocation, lessOrEqual bool) bool {
	if a.SameSize(other) {
		v := other.GetValue()
		condition := true
		for i := 0; i < len(a.x) && condition; i++ {
			if lessOrEqual {
				condition = a.x[i] <= v[i]
			} else {
				condition = a.x[i] == v[i]
			}
		}
		return condition
	}
	return false
}

// String : a print out of the allocation
func (a *Allocation) String() string {
	return fmt.Sprint(a.x)
}

// StringPretty : a print out of the allocation with resource names (empty if unequal lengths)
func (a *Allocation) StringPretty(resourceNames []string) string {
	n := len(a.x)
	if len(resourceNames) != n {
		return ""
	}
	var b bytes.Buffer
	b.WriteString("[")
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "%s:%d", resourceNames[i], a.x[i])
		if i < n-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString("]")
	return b.String()
}

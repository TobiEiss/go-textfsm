package process

import "reflect"

// Column of record
type Column struct {
	Type    reflect.Type
	Entries []interface{}
}

// String converts all entries to a slice of string
func (c Column) String() []string {
	strings := []string{}
	for _, entrie := range c.Entries {
		strings = append(strings, entrie.(string))
	}
	return strings
}

// Int converts entries to slice of int
func (c Column) Int() []int {
	ints := []int{}
	for _, entrie := range c.Entries {
		ints = append(ints, entrie.(int))
	}
	return ints
}

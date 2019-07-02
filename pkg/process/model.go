package process

import "reflect"

// Column of record
type Column struct {
	Type reflect.Type `json:"-"`
	Name string
}

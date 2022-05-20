package enum

import (
	"errors"
	"fmt"
	"togo/rest"
)

type enumThing int

const (
	_min enumThing = iota
	Thing1
	Thing2
	Thing3
	Thing4
	_max
)

var ErrNewThing = errors.New("invalid thing type")

type Thing struct {
	t enumThing
}

func NewThing(t enumThing) Thing {
	if t <= _min || t >= _max {
		rest.Panic(fmt.Errorf("%w: %v", ErrNewThing, t))
	}
	return Thing{t}
}

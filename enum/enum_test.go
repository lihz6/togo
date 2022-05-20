package enum_test

import (
	"errors"
	"testing"
	"togo/enum"
)

func TestNewThingMin(t *testing.T) {
	defer func() {
		err, ok := recover().(error)
		if ok && !errors.Is(err, enum.ErrNewThing) {
			t.Errorf("enum.NewThing(%v) should panic", enum.Min)
		}
	}()
	enum.NewThing(enum.Min)
}

func TestNewThingMax(t *testing.T) {
	defer func() {
		err, ok := recover().(error)
		if ok && !errors.Is(err, enum.ErrNewThing) {
			t.Errorf("enum.NewThing(%v) should panic", enum.Max)
		}
	}()
	enum.NewThing(enum.Max)
}

func TestNewThing(t *testing.T) {
	i := enum.Min + 1
	defer func() {
		err, ok := recover().(error)
		if ok && errors.Is(err, enum.ErrNewThing) {
			t.Errorf("enum.NewThing(%v) should not panic", i)
		}
	}()
	for ; i < enum.Max; i++ {
		enum.NewThing(i)
	}
}

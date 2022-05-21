package opts_test

import (
	"testing"
	"togo/opts"
)

func TestOpts(t *testing.T) {
	thing := opts.NewThing("required", opts.WithCache(true), opts.WithTitle("title"))
	if thing.SetCache(false) != true {
		t.Error("thing.SetCache(false)")
	}
	if thing.SetTitle("") != "title" {
		t.Error(`thing.SetTitle("")`)
	}
}

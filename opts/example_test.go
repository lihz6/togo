package opts_test

import "togo/opts"

func ExampleNewThing_args1() {
	opts.NewThing("required")
}

func ExampleNewThing_args2() {
	opts.NewThing("required", opts.WithCache(true))
}

func ExampleNewThing_args3() {
	opts.NewThing("required", opts.WithCache(true), opts.WithTitle("title"))
}

func ExampleThing_SetCache() {
	t := opts.NewThing("required")
	p := t.SetCache(false)
	defer t.SetCache(p)
}

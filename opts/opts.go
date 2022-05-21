package opts

type options struct {
	required string
	cache    bool
	title    string
}

type applyOption = func(opts *options)

func WithCache(cache bool) applyOption {
	return func(opts *options) {
		opts.cache = cache
	}
}

func WithTitle(title string) applyOption {
	return func(opts *options) {
		opts.title = title
	}
}

func NewThing(required string, applyOptions ...applyOption) Thing {
	opts := options{required: required}
	for _, applyOption := range applyOptions {
		applyOption(&opts)
	}
	return Thing{opts: opts}
}

type Thing struct {
	opts options
}

func (t *Thing) SetTitle(title string) string {
	t.opts.title, title = title, t.opts.title
	return title
}

func (t *Thing) SetCache(cache bool) bool {
	t.opts.cache, cache = cache, t.opts.cache
	return cache
}

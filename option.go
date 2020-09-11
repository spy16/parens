package parens

// Option for context type
type Option func(*Context)

// WithMapFactory sets the StackFrame implementation on the context.
func WithMapFactory(f MapFactory) Option {
	if f == nil {
		f = newBasicMap
	}

	return func(ctx *Context) {
		ctx.newMap = f
	}
}

// WithMaxRecursion sets the maximum allowable stack depth (including the global frame).
func WithMaxRecursion(depth uint) Option {
	if depth == 0 {
		panic("max recursion must be nonzero")
	}

	return func(ctx *Context) {
		ctx.maxDepth = int(depth)
	}
}

func withDefaults(opt []Option) []Option {
	return append([]Option{
		WithMapFactory(nil),
		WithMaxRecursion(10000),
	}, opt...)
}

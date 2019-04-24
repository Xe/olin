// Package opname contains an extensible "operation name" construct for go
// applications. This allows a user to create a base operation, eg:
// "createWidget" and then sub-operations such as "pgSaveWidget" and they will
// get composed such as "createWidget.pgSaveWidget" in log lines. Each
// operation name adds a lightweight "layer" to the context.
package opname

import (
	"context"
)

type ctxKey int

const key ctxKey = iota

// Get fetches the operation name from the given context.
func Get(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(key).(string)
	return val, ok
}

// With stores the current operation name to the context, optionally prepending
// the existing operation name in the context if it exists.
func With(ctx context.Context, name string) context.Context {
	prep, ok := Get(ctx)
	if ok {
		name = prep + "." + name
	}

	return context.WithValue(ctx, key, name)
}

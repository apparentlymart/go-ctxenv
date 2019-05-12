package ctxenv

import (
	"context"
	"os"
)

type ctxKeyType int

const ctxKey ctxKeyType = 1

// Environ returns a slice describing the environment, either from the given
// context or from the real process environment.
func Environ(ctx context.Context) []string {
	e, local := environ(ctx)
	if !local {
		return e // no need to copy, because this is already a copy from the syscall
	}

	// If it's local then we'll make a copy before we return, because otherwise
	// the caller would be able to mutate the supposedly-immutable environment
	// slice inside the context.
	if len(e) == 0 {
		return nil
	}
	ret := make([]string, len(e))
	copy(ret, e)
	return ret
}

// Getenv retrieves the value of the environment variable with the given name,
// either from the given context or from the real process environment.
func Getenv(ctx context.Context, key string) string {
	e, local := environCtxOnly(ctx)
	if !local {
		// Easy path, then
		return os.Getenv(key)
	}

	_, val := findInEnviron(e, key)
	return val
}

// WithEnviron creates a new context whose local environment is exactly as
// given, ignoring any existing entries in the parent context.
//
// This might be useful if you need to fully-control the environment variable
// table for a unit test, for example.
func WithEnviron(ctx context.Context, e []string) context.Context {
	// If we use the caller's slice then they would be able to mutate it,
	// so we'll copy it.
	ourE := make([]string, len(e))
	copy(ourE, e)
	return withEnviron(ctx, ourE)
}

// Setenv creates a new context whose local environment is the same as the
// given context except for overriding the given variable name with the given
// value.
//
// Each call to Setenv creates a copy of the parent environment table. The
// expected usage model here is that there will be only a few calls to this
// function early on during application startup and so the copying won't
// hurt too much. Calling Setenv in a tight loop is inadvisable.
func Setenv(ctx context.Context, key string, value string) context.Context {
	e, _ := environ(ctx)
	idx, _ := findInEnviron(e, key)
	if idx == -1 && value == "" {
		// If it wasn't already present and we were unsetting it anyway then
		// we'll just keep the same slice, since that's functionally equivalent
		// and saves some allocations and copying.
		return withEnviron(ctx, e)
	}

	if idx == -1 {
		// If there isn't already an entry for this key then we'll be appending
		// it to the end of our parent environ.
		ret := make([]string, len(e), len(e)+1)
		copy(ret, e)
		ret = append(ret, key+"="+value)
		return withEnviron(ctx, ret)
	}

	if value == "" {
		if len(e) == 1 { // will have nothing left after we remove this, then
			return Clearenv(ctx)
		}
		// If we're unsetting then we want to remove the entry entirely, not
		// overwrite it with an empty-string entry.
		ret := make([]string, len(e)-1, len(e)-1)
		copy(ret, e[:idx])
		copy(ret[idx:], e[idx+1:])
		return withEnviron(ctx, ret)
	}

	// If we get here then there _is_ an entry already present, so we're going
	// to overwrite it.
	ret := make([]string, len(e), len(e))
	copy(ret, e)
	ret[idx] = key + "=" + value
	return withEnviron(ctx, ret)
}

// Clearenv creates a new context whose local environment is totally empty.
func Clearenv(ctx context.Context) context.Context {
	return withEnviron(ctx, nil)
}

func environCtxOnly(ctx context.Context) ([]string, bool) {
	if got := ctx.Value(ctxKey); got != nil {
		return got.([]string), true
	}
	return nil, false

}

func environ(ctx context.Context) ([]string, bool) {
	e, local := environCtxOnly(ctx)
	if !local {
		e = os.Environ()
	}
	return e, local
}

func withEnviron(ctx context.Context, e []string) context.Context {
	return context.WithValue(ctx, ctxKey, e)
}

// findInEnviron returns the index of the first assignment of the given key
// in the given environ slice, or -1 if there is no such assignment.
func findInEnviron(e []string, key string) (int, string) {
	matchLen := len(key) + 1
	for i, ee := range e {
		if len(ee) < matchLen {
			continue // can't possibly match
		}
		if ee[:matchLen-1] == key && ee[matchLen-1] == '=' {
			return i, ee[matchLen:]
		}
	}
	return -1, ""
}

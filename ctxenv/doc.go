// Package ctxenv is an alternative to os.Getenv and os.Environ that allows
// environment variable values to optionally be locally overridden inside
// a context.Context.
//
// The goal is to allow software that normally accesses environment variables
// directly to be given localized environment variable values in unusual
// situations where that is necessary, such as when writing tests that can
// be run concurrently.
//
// This is not a fully-general mechanism: it works only when both the caller
// and the callee agree to use this mechanism. Therefore it is appropriate
// only when caller and callee are somewhat coupled, such as when the caller
// is a unit test for the callee.
//
// However, it's designed to fall back automatically to the "real" environment
// variables when the caller does not participate in ctxenv, allowing functions
// to make use of this mechanism without imposing any new requirements on the
// caller in the common case where normal process-level environment variables
// are enough.
package ctxenv

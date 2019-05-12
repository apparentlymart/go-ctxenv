# ctxenv

[![godoc reference](https://godoc.org/github.com/apparentlymart/go-ctxenv/ctxenv?status.svg)](https://godoc.org/github.com/apparentlymart/go-ctxenv/ctxenv)

`ctxenv` is a small Go package that allows context-scoped environment variable
overrides for co-operating functions.

Sometimes a function will directly access environment variables to customize
its behavior. While passing in the settings from the caller would be more
flexible, we might choose to access the environment directly in situations
where the setting is primarily aimed at the end-user and the calling function
needs to override it only for some unusual reason, such as unit testing with
different permutations.

The functions in `ctxenv` mimic the functions of the same name in the `os`
package, but rather than interacting directly with the process environment
via system calls they first check to see if the given context contains an
overridden context-local environment table, preferring that if present.

Due to the context-handling behavior, these functions have a slightly different
usage pattern than their `os` equivalents:

- All of them take `ctx context.Context` as an additional first argument.
- The functions that would normally mutate the environment table instead
  return a new context containing their result as an override.
- If given a context with no overridden environment table, the reading
  functions will read from the real environment instead, so the caller doesn't
  need to worry about whether an override table is present.

Overridden environment tables are immutable, but the "real" environment is not.
Each override table is a full copy of the parent with some modifications
applied, so the first override table effectively freezes the set of values
in the environment when it is created, and will not reflect any mutations of
the real environment made by calls to `os.Setenv` or `os.Clearenv`.

## Using `ctxenv` in your application

If you are writing a function that customizes its behavior based on the
environment, you can use `ctxenv` by ensuring that the function accepts a
context as an argument and by using `ctxenv.Getenv` instead of `os.Getenv`.

If you do so, you may wish to include in your documentation that the function
supports environment variable overrides via this package, so that callers will
know that they can then use `ctx.Setenv` and `ctx.Clearenv` to locally override
the main environment table when needed.

---

This package is considered feature-complete and "done", so no additional
features will be added.

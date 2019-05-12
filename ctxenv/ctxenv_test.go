package ctxenv

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Most of our tests access the real process environment at some point,
	// so we'll make sure its contents are predictable regardless of where
	// we are running.
	os.Clearenv()
	os.Setenv("CTXENV_EXAMPLE", "real environment")
	os.Exit(m.Run())
}

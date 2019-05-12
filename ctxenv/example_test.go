package ctxenv_test

import (
	"context"
	"fmt"

	"github.com/apparentlymart/go-ctxenv/ctxenv"
)

func Example() {
	rootCtx := context.Background()
	fmt.Printf("root has %q\n\n", ctxenv.Getenv(rootCtx, "CTXENV_EXAMPLE"))

	overrideCtx := ctxenv.Setenv(rootCtx, "CTXENV_EXAMPLE", "locally overridden")
	fmt.Printf("override has %q\n", ctxenv.Getenv(overrideCtx, "CTXENV_EXAMPLE"))
	fmt.Printf("root still has %q\n\n", ctxenv.Getenv(rootCtx, "CTXENV_EXAMPLE"))

	clearCtx := ctxenv.Clearenv(overrideCtx)
	fmt.Printf("clear has %q\n", ctxenv.Getenv(clearCtx, "CTXENV_EXAMPLE"))
	fmt.Printf("override still has %q\n", ctxenv.Getenv(overrideCtx, "CTXENV_EXAMPLE"))
	fmt.Printf("root still has %q\n\n", ctxenv.Getenv(rootCtx, "CTXENV_EXAMPLE"))

	fmt.Println("rootCtx environment:")
	for _, ee := range ctxenv.Environ(rootCtx) {
		fmt.Printf("- %s\n", ee)
	}

	fmt.Println("\noverrideCtx environment:")
	for _, ee := range ctxenv.Environ(overrideCtx) {
		fmt.Printf("- %s\n", ee)
	}

	// Output:
	// root has "real environment"
	//
	// override has "locally overridden"
	// root still has "real environment"
	//
	// clear has ""
	// override still has "locally overridden"
	// root still has "real environment"
	//
	// rootCtx environment:
	// - CTXENV_EXAMPLE=real environment
	//
	// overrideCtx environment:
	// - CTXENV_EXAMPLE=locally overridden
}

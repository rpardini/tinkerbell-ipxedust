package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/equinix-labs/otel-init-go/otelinit"
	"github.com/rpardini/tinkerbell-ipxedust"
)

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	defer done()
	ctx, otelShutdown := otelinit.InitOpenTelemetry(ctx, "github.com/rpardini/tinkerbell-ipxedust")
	defer otelShutdown(ctx)

	if err := ipxedust.Execute(ctx, os.Args[1:]); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
	}
}

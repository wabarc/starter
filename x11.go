package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func runXvfb(ctx context.Context) error {
	// Xvfb $DISPLAY -ac -screen 0 $XVFB_WHD +extension GLX +render -noreset -nolisten tcp > /dev/null 2>&1 &
	display := os.Getenv("DISPLAY")
	whd := os.Getenv("XVFB_WHD")
	opts := []string{
		display,
		"-ac",
		"-screen",
		"0",
		whd,
		"-nolisten",
		"tcp",
	}
	cmd := exec.CommandContext(ctx, "Xvfb", opts...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Run StdoutPipe failed: %w", err)
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Start failed: %w", err)
	}

	go readOutput(out)
	go func() {
		_ = cmd.Wait()
	}()

	return nil
}

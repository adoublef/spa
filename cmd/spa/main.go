package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
)

type command interface {
	parse(args []string, getenv func(string) string) error
	run(ctx context.Context) error
}

var cmds = map[string]command{
	"serve": cmdServe,
}

func main() {
	err := run(context.Background(), os.Args[1:], os.Getenv)
	if errors.Is(err, flag.ErrHelp) {
		os.Exit(2)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "ERRO: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, getenv func(string) string) error {
	var cmd string
	if len(args) > 0 {
		cmd, args = args[0], args[1:]
	}
	for name, c := range cmds {
		// todo: handle help command
		// https://github.com/superfly/litefs/blob/main/cmd/litefs/main.go#L72C1-L77C55
		if name != cmd {
			continue
		}
		err := c.parse(args, getenv)
		if err != nil {
			return fmt.Errorf("%v: %w", err, flag.ErrHelp)
		}
		return c.run(ctx)
	}
	return fmt.Errorf("unknown command: %s: %w", cmd, flag.ErrHelp)
}

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
	"tmp.adoublef/spa/internal/net/http"
)

var cmdServe = &serve{}

type serve struct {
	Addr string
	Path http.Dir
}

func (c *serve) parse(args []string, _ func(string) string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.StringVar(&c.Addr, "addr", ":80", "http listening port")
	fs.Usage = func() {
		fmt.Println(`
The export command will start the server.

Usage:

	serve [arguments] PATH

Arguments:
`[1:])
		fs.PrintDefaults()
		fmt.Println("")
	}
	if err := fs.Parse(args); err != nil {
		return err
	} else if fs.NArg() == 0 {
		fs.Usage()
		return flag.ErrHelp
	} else if fs.NArg() > 1 {
		return fmt.Errorf("too many arguments")
	}
	// i could use parse another flag but a positional arg is cool
	c.Path = http.Dir(fs.Arg(0))
	return nil
}

func (c *serve) run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	s := &http.Server{
		Addr:        c.Addr,
		Handler:     http.Handler(c.Path),
		BaseContext: func(l net.Listener) context.Context { return ctx },
		// todo: timeouts
	}
	s.RegisterOnShutdown(cancel)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() (err error) {
		switch {
		case s.TLSConfig != nil:
			err = s.ListenAndServeTLS("", "")
		default:
			err = s.ListenAndServe()
		}
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.Shutdown(ctx)
		if err != nil {
			if ec := s.Close(); ec != nil {
				err = ec
			}
		}
		return err
	})

	// no reason for this check since there is nothing beyond this stage
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// printUsage prints the help screen to STDOUT.
func printUsage() {
	fmt.Println(`
a.out is a single page application hosted on a simple web server.

Usage:

	a.out <command> [arguments]

The commands are:

	serve       start web server
`[1:])
}

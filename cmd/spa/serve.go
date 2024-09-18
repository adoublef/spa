package main

import (
	"context"
	"errors"
	"flag"
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
	// todo: static files
}

func (c *serve) parse(args []string, _ func(string) string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.StringVar(&c.Addr, "addr", ":80", "http listening port")
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	return nil
}

func (c *serve) run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	s := &http.Server{
		Addr:        c.Addr,
		Handler:     http.Handler(),
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

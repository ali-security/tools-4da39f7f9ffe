// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"
)

const (
	ExampleLine   = exampleLine
	ExampleColumn = exampleColumn
	ExampleOffset = exampleOffset
)

const implicitlyBindsAll = "implicitly binds all network interfaces"

// runServe runs the serve command in a goroutine, so that a server that
// actually starts listening fails the test instead of hanging it.
func runServe(t *testing.T, s *Serve) error {
	t.Helper()
	errc := make(chan error, 1)
	go func() { errc <- s.Run(context.Background()) }()
	select {
	case err := <-errc:
		return err
	case <-time.After(10 * time.Second):
		t.Fatalf("serve -listen=%s did not return: the server is listening", s.Address)
		return nil
	}
}

func TestServeListenRejectsImplicitAllInterfaces(t *testing.T) {
	for _, addr := range []string{":0", ":"} {
		s := &Serve{Address: addr, app: &Application{}}
		err := runServe(t, s)
		if err == nil || !strings.Contains(err.Error(), implicitlyBindsAll) {
			t.Errorf("serve -listen=%s: got error %v, want %q", addr, err, implicitlyBindsAll)
		}
	}
}

func TestServeListenAllowsExplicitHost(t *testing.T) {
	// The port is out of range, so net.Listen fails without binding anything;
	// the address must reach it rather than be rejected by the -listen check.
	for _, addr := range []string{"localhost:99999", "127.0.0.1:99999"} {
		s := &Serve{Address: addr, app: &Application{}}
		err := runServe(t, s)
		if err == nil || strings.Contains(err.Error(), implicitlyBindsAll) {
			t.Errorf("serve -listen=%s: got error %v, want a listen error", addr, err)
		}
	}
}

func TestServePortFlagRemoved(t *testing.T) {
	// -port=N listened on ":N", i.e. on all network interfaces.
	typ := reflect.TypeOf(Serve{})
	for i := 0; i < typ.NumField(); i++ {
		if f := typ.Field(i); f.Tag.Get("flag") == "port" {
			t.Errorf("Serve.%s still registers the -port flag", f.Name)
		}
	}
}

// Copyright 2017 The 99c Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/cznic/virtual"
)

func exit(code int, msg string, arg ...interface{}) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, os.Args[0]+": "+msg, arg...)
	}
	os.Exit(code)
}

func main() {
	if !profile {
		exit(1, `This tool must be built using '-tags virtual.profile', please rebuild it.
Use
	$ go install -tags virtual.profile github.com/cznic/99c/99prof
or
	$ go get [-u] -tags virtual.profile github.com/cznic/99c/99prof
`)
	}

	functions := flag.Bool("functions", false, "")
	instructions := flag.Bool("instructions", false, "")
	lines := flag.Bool("lines", false, "")
	rate := flag.Int("rate", 1000, "")
	flag.Parse()

	if flag.NArg() == 0 {
		exit(2, "missing program name %v\n", os.Args)
	}

	nm := flag.Arg(0)
	bin, err := os.Open(nm)
	if err != nil {
		exit(1, "%v\n", err)
	}

	var b virtual.Binary
	if _, err := b.ReadFrom(bufio.NewReader(bin)); err != nil {
		exit(1, "%v\n", err)
	}

	args := os.Args[1:]
	for i, v := range args {
		if v == nm {
			args = args[i+1:]
			break
		}
	}

	var opts []virtual.Option
	if *functions {
		opts = append(opts, virtual.ProfileFunctions())
	}
	if *lines {
		opts = append(opts, virtual.ProfileLines())
	}
	if *instructions {
		opts = append(opts, virtual.ProfileInstructions())
	}
	if n := *rate; n != 0 {
		opts = append(opts, virtual.ProfileRate(n))
	}
	t0 := time.Now()

	code, err := virtual.Exec(&b, os.Args[1:], os.Stdin, os.Stdout, os.Stderr, 0, 8<<20, "", opts...)
	if err != nil {
		if code == 0 {
			code = 1
		}
		exit(code, "%v\n", err)
	}

	d := time.Since(t0)
	_ = d
	panic("TODO")

	exit(code, "")
}

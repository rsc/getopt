// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package getopt

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
	"testing"
)

// TODO: Real tests.

type testFlagSet struct {
	A    *bool
	B    *bool
	C    *bool
	D    *bool
	I    *int
	Long *int
	S    *string
	Args []string

	t    *testing.T
	buf  bytes.Buffer
	flag *FlagSet
}

func (tf *testFlagSet) str() string {
	out := ""
	if *tf.A {
		out += " -a"
	}
	if *tf.B {
		out += " -b"
	}
	if *tf.C {
		out += " -c"
	}
	if *tf.D {
		out += " -d"
	}
	if *tf.I != 0 {
		out += fmt.Sprintf(" -i %d", *tf.I)
	}
	if *tf.Long != 0 {
		out += fmt.Sprintf(" --long %d", *tf.Long)
	}
	if *tf.S != "" {
		out += " -s " + *tf.S
	}
	if len(tf.Args) > 0 {
		out += " " + strings.Join(tf.Args, " ")
	}
	if out == "" {
		return out
	}
	return out[1:]
}

func newTestFlagSet(t *testing.T) *testFlagSet {
	tf := &testFlagSet{t: t, flag: NewFlagSet("x", flag.ContinueOnError)}
	f := tf.flag
	f.SetOutput(&tf.buf)
	tf.A = f.Bool("a", false, "desc of a")
	tf.B = f.Bool("b", false, "desc of b")
	tf.C = f.Bool("c", false, "desc of c")
	tf.D = f.Bool("d", false, "desc of d")
	tf.Long = f.Int("long", 0, "long only")
	f.Alias("a", "aah")
	f.Aliases("b", "beeta", "c", "charlie")
	tf.I = f.Int("i", 0, "i")
	f.Alias("i", "india")
	tf.S = f.String("sierra", "", "string")
	f.Alias("s", "sierra")

	return tf
}

var tests = []struct {
	cmd string
	out string
}{
	{"-i 1", "-i 1"},
	{"--india 1", "-i 1"},
	{"--india=1", "-i 1"},
	{"-i=1", `ERR: invalid value "=1" for flag -i: strconv.ParseInt: parsing "=1": invalid syntax`},
	{"--i=1", "-i 1"},
	{"-abc", "-a -b -c"},
	{"--abc", `ERR: flag provided but not defined: --abc`},
	{"-sfoo", "-s foo"},
	{"-s foo", "-s foo"},
	{"--s=foo", "-s foo"},
	{"-s=foo", "-s =foo"},
	{"-s", `ERR: missing argument for -s`},
	{"--s", `ERR: missing argument for --s`},
	{"--s=", ``},
	{"-sfooi1 -i2", "-i 2 -s fooi1"},
	{"-absfoo", "-a -b -s foo"},
	{"-i1 -- arg", "-i 1 arg"},
	{"-i1 - arg", "-i 1 - arg"},
	{"-i1 --- arg", `ERR: flag provided but not defined: ---`},
	{"-i1 arg", "-i 1 arg"},
	{"--aah --charlie --beeta --sierra=123", "-a -b -c -s 123"},
	{"-i1 --long=2", "-i 1 --long 2"},
}

func TestBasic(t *testing.T) {
	for _, tt := range tests {
		tf := newTestFlagSet(t)
		err := tf.flag.Parse(strings.Fields(tt.cmd))
		var out string
		if err != nil {
			out = "ERR: " + err.Error()
		} else {
			tf.Args = tf.flag.Args()
			out = tf.str()
		}
		if out != tt.out {
			t.Errorf("%s:\nhave %s\nwant %s", tt.cmd, out, tt.out)
		}
	}
}

var wantHelpText = `  -a, --aah
    	desc of a
  -b, --beeta
    	desc of b
  -c, --charlie
    	desc of c
  -d	desc of d
  -i, --india int
    	i
  --long int
    	long only
  -s, --sierra string
    	string
`

func TestHelpText(t *testing.T) {
	tf := newTestFlagSet(t)
	tf.flag.PrintDefaults()
	out := tf.buf.String()
	if out != wantHelpText {
		t.Errorf("have<\n%s>\nwant<\n%s>", out, wantHelpText)
	}
}

# rsc.io/getopt

[For full package documentation, see [https://godoc.org/rsc.io/getopt](https://godoc.org/rsc.io/getopt).]

    package getopt // import "rsc.io/getopt"

Package getopt parses command lines using [_getopt_(3)](http://man7.org/linux/man-pages/man3/getopt.3.html) syntax. It is a
replacement for `flag.Parse` but still expects flags themselves to be defined
in package flag.

Flags defined with one-letter names are available as short flags (invoked
using one dash, as in `-x`) and all flags are available as long flags (invoked
using two dashes, as in `--x` or `--xylophone`).

To use, define flags as usual with [package flag](https://godoc.org/flag). Then introduce any aliases
by calling `getopt.Alias`:

    getopt.Alias("v", "verbose")

Or call `getopt.Aliases` to define a list of aliases:

    getopt.Aliases(
    	"v", "verbose",
    	"x", "xylophone",
    )

One name in each pair must already be defined in package flag (so either
"v" or "verbose", and also either "x" or "xylophone").

Then parse the command-line:

    getopt.Parse()

If it encounters an error, `Parse` calls `flag.Usage` and then exits the
program.

When writing a custom `flag.Usage` function, call `getopt.PrintDefaults` instead
of `flag.PrintDefaults` to get a usage message that includes the
names of aliases in flag descriptions.

At initialization time, package getopt installs a new `flag.Usage` that is the same
as the default `flag.Usage` except that it calls `getopt.PrintDefaults` instead
of `flag.PrintDefaults`.

This package also defines a `FlagSet` wrapping the standard `flag.FlagSet`.

## Caveat

In general Go flag parsing is preferred for new programs, because it is not
as pedantic about the number of dashes used to invoke a flag (you can write
`-verbose` or `--verbose` and the program does not care). This package is meant
to be used in situations where, for legacy reasons, it is important to use
exactly _getopt_(3) syntax, such as when rewriting in Go an existing tool that
already uses _getopt_(3).

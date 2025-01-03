package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"

	"github.com/a-h/templ/cmd/templ/sloghandler"
	"github.com/fatih/color"
	"github.com/garrettladley/templ-static/cmd/templ-static/generatecmd"
)

func main() {
	code := run(os.Stdin, os.Stdout, os.Stderr, os.Args)
	if code != 0 {
		os.Exit(code)
	}
}

const usageText = `usage: templ-static <command> [<args>...]

templ-static - <TODO>

See docs at <TODO>

commands:
  <TODO>   <TODO>
`

func run(stdin io.Reader, stdout, stderr io.Writer, args []string) (code int) {
	if len(args) < 2 {
		fmt.Fprint(stderr, usageText)
		return 64 // EX_USAGE
	}
	switch args[1] {
	case "generate":
		return generateCmd(stdout, stderr, args[2:])
	case "help", "-help", "--help", "-h":
		fmt.Fprint(stdout, usageText)
		return 0
	}
	fmt.Fprint(stderr, usageText)
	return 64 // EX_USAGE
}

func newLogger(logLevel string, verbose bool, stderr io.Writer) *slog.Logger {
	if verbose {
		logLevel = "debug"
	}
	level := slog.LevelInfo.Level()
	switch logLevel {
	case "debug":
		level = slog.LevelDebug.Level()
	case "warn":
		level = slog.LevelWarn.Level()
	case "error":
		level = slog.LevelError.Level()
	}
	return slog.New(sloghandler.NewHandler(stderr, &slog.HandlerOptions{
		AddSource: logLevel == "debug",
		Level:     level,
	}))
}

const generateUsageText = `usage: templ-static generate [<args>...]

<TODO>

Args:
  -path <path>
    Generates code for all files in path. (default .)
  -v
    Set log verbosity level to "debug". (default "info")
  -log-level
    Set log verbosity level. (default "info", options: "debug", "info", "warn", "error")
  -help
    Print help and exit.

Examples:

  <TODO>:

    templ-static generate
`

func generateCmd(stdout, stderr io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("generate", flag.ExitOnError)
	pathFlag := cmd.String("path", ".", "")
	verboseFlag := cmd.Bool("v", false, "")
	logLevelFlag := cmd.String("log-level", "info", "")
	helpFlag := cmd.Bool("help", false, "")
	err := cmd.Parse(args)
	if err != nil {
		fmt.Fprint(stderr, generateUsageText)
		return 64 // EX_USAGE
	}
	if *helpFlag {
		fmt.Fprint(stdout, generateUsageText)
		return
	}

	log := newLogger(*logLevelFlag, *verboseFlag, stderr)

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Fprintln(stderr, "Stopping...")
		cancel()
	}()

	err = generatecmd.Run(ctx, log, generatecmd.Arguments{
		Path: *pathFlag,
	})
	if err != nil {
		color.New(color.FgRed).Fprint(stderr, "(✗) ")
		fmt.Fprintln(stderr, "Command failed: "+err.Error())
		return 1
	}
	return 0
}

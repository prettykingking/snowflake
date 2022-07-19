package cmd

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prettykingking/snowflake/pkg/config"
	"github.com/prettykingking/snowflake/pkg/util"
)

const (
	ExitCodeFailedStartup = iota
)

// Command represents a subcommand. Name, Func,
// and Short are required.
type Command struct {
	// The name of the subcommand. Must conform to the
	// format described by the RegisterCommand() godoc.
	// Required.
	Name string

	// The static configuration can be read from file
	// or environment variable
	Configuration string

	// Func is a function that executes a subcommand using
	// the parsed flags. It returns an exit code and any
	// associated error.
	// Required.
	Func CommandFunc

	// Usage is a brief message describing the syntax of
	// the subcommand's flags and args. Use [] to indicate
	// optional parameters and <> to enclose literal values
	// intended to be replaced by the user. Do not prefix
	// the string with "caddy" or the name of the command
	// since these will be prepended for you; only include
	// the actual parameters for this command.
	Usage string

	// Short is a one-line message explaining what the
	// command does. Should not end with punctuation.
	// Required.
	Short string

	// Long is the full help text shown to the user.
	// Will be trimmed of whitespace on both ends before
	// being printed.
	Long string

	// Flags is the flagset for command.
	Flags *flag.FlagSet
}

// Flags wraps a FlagSet so that typed values
// from flags can be easily retrieved.
type Flags struct {
	*flag.FlagSet
}

// String returns the string representation of the
// flag given by name. It panics if the flag is not
// in the flag set.
func (f Flags) String(name string) string {
	return f.FlagSet.Lookup(name).Value.String()
}

// Bool returns the boolean representation of the
// flag given by name. It returns false if the flag
// is not a boolean type. It panics if the flag is
// not in the flag set.
func (f Flags) Bool(name string) bool {
	val, _ := strconv.ParseBool(f.String(name))
	return val
}

// Int returns the integer representation of the
// flag given by name. It returns 0 if the flag
// is not an integer type. It panics if the flag is
// not in the flag set.
func (f Flags) Int(name string) int {
	val, _ := strconv.ParseInt(f.String(name), 0, strconv.IntSize)
	return int(val)
}

// Float64 returns the float64 representation of the
// flag given by name. It returns false if the flag
// is not a float64 type. It panics if the flag is
// not in the flag set.
func (f Flags) Float64(name string) float64 {
	val, _ := strconv.ParseFloat(f.String(name), 64)
	return val
}

// Duration returns the duration representation of the
// flag given by name. It returns false if the flag
// is not a duration type. It panics if the flag is
// not in the flag set.
func (f Flags) Duration(name string) time.Duration {
	val, _ := util.ParseDuration(f.String(name))
	return val
}

// CommandFunc is a command's function. It runs the
// command and returns the proper exit code along with
// any error that occurred.
type CommandFunc func(Flags, *config.Configuration) (int, error)

var commandNameRegex = regexp.MustCompile(`^[a-z0-9]$|^([a-z0-9]+-?[a-z0-9]*)+[a-z0-9]$`)
var commands = make(map[string]Command)

func init() {
	RegisterCommand(Command{
		Name:  "help",
		Func:  cmdHelp,
		Usage: "<command>",
		Short: "Show help for command",
	})
}

func Execute() {
	switch len(os.Args) {
	case 0:
		fmt.Printf("[FATAL] no arguments provided by OS; args[0] must be command\n")
		os.Exit(ExitCodeFailedStartup)
	case 1:
		os.Args = append(os.Args, "help")
	}

	subName := os.Args[1]
	sub, ok := commands[subName]
	if !ok {
		if strings.HasPrefix(subName, "-") {
			fmt.Println("[ERROR] first argument must be a subcommand; see 'help'")
		} else {
			fmt.Printf("[ERROR] '%s' is not a recognized subcommand; see 'snowflake help'\n", os.Args[1])
		}
		os.Exit(ExitCodeFailedStartup)
	}

	flSub := sub.Flags
	if flSub == nil {
		flSub = flag.NewFlagSet(sub.Name, flag.ExitOnError)
	}

	err := flSub.Parse(os.Args[2:])
	if err != nil {
		fmt.Println(err)
		os.Exit(ExitCodeFailedStartup)
	}

	var flSuper = Flags{flSub}
	cf := config.NewConfiguration()
	if flSuper.Lookup("config") != nil {
		ok, err = config.LoadFile(flSuper.String("config"), cf)
		if !ok {
			fmt.Printf("error: failed to parse configuration file. %s\n", err)
			os.Exit(ExitCodeFailedStartup)
		}
	}

	exitCode, err := sub.Func(flSuper, cf)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", sub.Name, err)
	}

	os.Exit(exitCode)
}

// RegisterCommand registers the command cmd.
// cmd.Name must be unique and conform to the
// following format:
//
//    - lowercase
//    - alphanumeric and hyphen characters only
//    - cannot start or end with a hyphen
//    - hyphen cannot be adjacent to another hyphen
//
// This function panics if the name is already registered,
// if the name does not meet the described format, or if
// any of the fields are missing from cmd.
//
// This function should be used in init().
func RegisterCommand(cmd Command) {
	if cmd.Name == "" {
		panic("command name is required")
	}
	if cmd.Func == nil {
		panic("command function missing")
	}
	if cmd.Short == "" {
		panic("command short string is required")
	}
	if _, exists := commands[cmd.Name]; exists {
		panic("command already registered: " + cmd.Name)
	}
	if !commandNameRegex.MatchString(cmd.Name) {
		panic("invalid command name")
	}
	commands[cmd.Name] = cmd
}

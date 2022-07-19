package cmd

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/prettykingking/snowflake/pkg/config"
)

func cmdHelp(fl Flags, _ *config.Configuration) (int, error) {
	args := fl.Args()

	if len(args) == 0 {
		desc := `Settings - Generate Globally Unique Identifier

Usage: snowflake <command> [<arg...>]

Commands:
`
		keys := make([]string, 0, len(commands))
		for k := range commands {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		for _, k := range keys {
			sub := commands[k]
			short := sub.Short
			desc += fmt.Sprintf("  %-15s %s\n", sub.Name, short)
		}

		desc += "\nUse snowflake help <command> for more information about a command.\n"

		fmt.Print(desc)

		return 0, nil
	} else if len(args) > 1 {
		return 1, errors.New("can only give help with one command")
	}

	sub, ok := commands[args[0]]
	if !ok {
		return 1, fmt.Errorf("unkown command %s", args[0])
	}

	helpText := strings.TrimSpace(sub.Long)
	result := fmt.Sprintf("%s\n\nUsage: snowflake %s %s\n",
		helpText, sub.Name, strings.TrimSpace(sub.Usage),
	)

	if help := flagHelp(sub.Flags); help != "" {
		result += fmt.Sprintf("\nOptions:\n%s", help)
	}

	fmt.Print(result)

	return 0, nil
}

// flagHelp returns the help text for fs.
func flagHelp(fl *flag.FlagSet) string {
	if fl == nil {
		return ""
	}

	// temporarily redirect output
	out := fl.Output()
	defer fl.SetOutput(out)

	buf := new(bytes.Buffer)
	fl.SetOutput(buf)
	fl.PrintDefaults()
	return buf.String()
}

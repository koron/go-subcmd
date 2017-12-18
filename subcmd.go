package subcmd

import (
	"flag"
	"fmt"
	"sort"
	"strings"
)

// Main represents entry point of a sub-command
type Main func([]string) error

// Main2 represents entry point of a sub-command
type Main2 func(*flag.FlagSet, []string) error

// Subcmds represents sub-commands table. This should has Main, Main2 or
// Subcmds as values.
type Subcmds map[string]interface{}

// Run parses args and executes one of sub-commands.
func (sc Subcmds) Run(args []string) error {
	return sc.RunWithName("", args)
}

// RunWithName parses args and executes one of sub-commands with name of
// flag.FlagSet which passed to Main2 entry points.
func (sc Subcmds) RunWithName(name string, args []string) error {
	// nested subcmds is not so deep. it would be enough for 8.
	nc := len(args) + 1
	if nc > 8 {
		nc = 8
	}
	cmds := make([]string, 0, len(args)+1)
	if name != "" {
		cmds = append(cmds, name)
	}
	return sc.run(cmds, args)
}

func (sc Subcmds) run(cmds, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("require one of sub-commands: %s", sc.names())
	}
	n, remains := args[0], args[1:]
	v, ok := sc[n]
	if !ok {
		return fmt.Errorf("unknown %q sub-command, should be one of: %s", n, sc.names())
	}
	cmds = append(cmds, n)
	switch w := v.(type) {
	case Main:
		return w(remains)
	case Main2:
		return sc.kickMain2(w, cmds, remains)
	case Subcmds:
		return w.run(cmds, remains)
	default:
		return fmt.Errorf("unexpected subcmds: %s", n)
	}
}

func (sc Subcmds) names() string {
	a := make([]string, 0, len(sc))
	for k := range sc {
		a = append(a, k)
	}
	sort.Strings(a)
	return strings.Join(a, ", ")
}

func (sc Subcmds) kickMain2(m Main2, cmds, args []string) error {
	n := strings.Join(cmds, " ")
	fs := flag.NewFlagSet(n, flag.ExitOnError)
	err := m(fs, args)
	if err != nil {
		return err
	}
	return nil
}

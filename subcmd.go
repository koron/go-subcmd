package subcmd

import (
	"fmt"
	"sort"
	"strings"
)

type Main func([]string) error

type Subcmds map[string]interface{}

func (sc Subcmds) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("require one of sub-commands: %s", sc.names())
	}
	n, remains := args[0], args[1:]
	v, ok := sc[n]
	if !ok {
		return fmt.Errorf("unknown %q sub-command, should be one of: %s", n, sc.names())
	}
	switch w := v.(type) {
	case Main:
		return w(remains)
	case Subcmds:
		return w.Run(remains)
	default:
		return fmt.Errorf("unexpected subcmds: %s", n)
	}
}

func (sc Subcmds) names() string {
	a := make([]string, 0, len(sc))
	for k, _ := range sc {
		a = append(a, k)
	}
	sort.Strings(a)
	return strings.Join(a, ", ")
}

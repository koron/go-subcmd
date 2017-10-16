package subcmd

import (
	"bytes"
	"flag"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var cmds = Subcmds{
	"foo": Subcmds{
		"1": Main(foo1),
		"2": Main2(foo2),
	},
}

var (
	cmdsFunc   string
	cmdsFS     *flag.FlagSet
	cmdsFSName string
	cmdsArgs   []string
)

func foo1(args []string) error {
	cmdsFunc = "foo1"
	cmdsFS = nil
	cmdsFSName = ""
	cmdsArgs = args
	return nil
}

func foo2(fs *flag.FlagSet, args []string) error {
	cmdsFunc = "foo2"
	cmdsFS = fs
	n, err := name(fs)
	if err != nil {
		return err
	}
	cmdsFSName = n
	cmdsArgs = args
	return nil
}

func name(fs *flag.FlagSet) (string, error) {
	bb := &bytes.Buffer{}
	fs.SetOutput(bb)
	fs.Usage()
	s, err := bb.ReadString('\n')
	if err != nil {
		return "", err
	}
	if s == "Usage:\n" {
		return "", nil
	}
	const (
		prefix = "Usage of "
		suffix = ":\n"
	)
	if !strings.HasPrefix(s, prefix) || !strings.HasSuffix(s, suffix) {
		return "", fmt.Errorf("unexpected prefix or suffix: %s", s)
	}
	return s[len(prefix) : len(s)-len(suffix)], nil
}

type Expect struct {
	Func   string
	Args   []string
	HasFS  bool
	FSName string
}

func (ex *Expect) Check() error {
	if cmdsFunc != ex.Func {
		return fmt.Errorf("func doesn't match: expected=%q actual=%q", ex.Func, cmdsFunc)
	}
	if (cmdsFS != nil) != ex.HasFS {
		return fmt.Errorf("providing FlagSet doesn't match: expected=%t actual=%t", ex.HasFS, cmdsFS != nil)
	}
	if ex.HasFS {
		if cmdsFSName != ex.FSName {
			return fmt.Errorf("name of FlagSet doesn't match: expected=%q actual=%q", ex.FSName, cmdsFSName)
		}
	}
	if !reflect.DeepEqual(cmdsArgs, ex.Args) {
		return fmt.Errorf("args doesn't match: expected=%+v actual=%+v", ex.Args, cmdsArgs)
	}
	return nil
}

func TestCmds(t *testing.T) {
	ok := func(args []string, exp *Expect, msg string) {
		err := cmds.RunWithName("TestCmds", args)
		if err != nil {
			t.Fatalf("%s: RunWithName() failed: %s", msg, err)
		}
		err = exp.Check()
		if err != nil {
			t.Errorf("%s: unexpected: %s", msg, err)
		}
	}

	ok([]string{"foo", "1", "aaa", "bbb"}, &Expect{
		Func: "foo1", Args: []string{"aaa", "bbb"},
		HasFS: false,
	}, "run foo1")
	ok([]string{"foo", "2", "aaa", "bbb"}, &Expect{
		Func: "foo2", Args: []string{"aaa", "bbb"},
		HasFS: true, FSName: "TestCmds foo 2",
	}, "run foo1")
}

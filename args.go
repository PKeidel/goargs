package goargs

import (
	"fmt"
	"os"
	"strings"
)

// [x] Support strings
// [ ] Support ints
// [ ] Support flags (bool)
// [x] Support os.Args parsing
// [x] Support env parsing
// [ ] Support config file parsing
// [ ] Support multi flag to allow arrays like `--host host01 host02`
// [ ] Detect missing args
// [ ] Detect undefined args
// [ ] Detect wrong types
// [ ] Create go generator for Arg struct

type ArgType uint8
type SourceType uint8

//go:generate stringer -type=ArgType
const (
	TypeString ArgType = iota
	TypeBool
	TypeInt
)

//go:generate stringer -type=SourceType
const (
	SourceNotFound SourceType = iota
	SourceFile
	SourceEnv
	SourceArgs
)

var (
	DefaultArgs []Arg
	DefaultSourceOrder = []SourceType{SourceFile, SourceEnv, SourceArgs}
	EnvPrefix = "GOARG_"

	values map[string]interface{} = make(map[string]interface{}, 0)
	sourceNotFound SourceType = SourceNotFound
)

type Arg struct {
	argtype ArgType
	longName, shortName string
	multi bool
	source *SourceType
	refS *string
}

func (a *Arg) String() string {
	return fmt.Sprintf(
		"Arg[type: %s, longName: %s, shortName: %s, source: %s, ref: %s]",
		a.argtype,
		a.longName,
		a.shortName,
		a.source,
		*a.refS,
	)
}

func WithStringS(longName, shortName string, ref *string) {
	DefaultArgs = append(DefaultArgs, Arg{
		argtype: TypeString,
		longName: longName,
		shortName: shortName,
		multi: false,
		source: &sourceNotFound,
		refS: ref,
	})
}

func WithString(longName string, ref *string) {
	WithStringS(longName, "", ref)
}

func ParseWith(args, envs []string) {
	for _, source := range DefaultSourceOrder {
		switch(source) {
		case SourceFile:
			parseFile()
		case SourceEnv:
			parseEnv(envs)
		case SourceArgs:
			parseArgs(args)
		}
	}

	log("Parse() result:")
	for _, a := range DefaultArgs {
		logf("    %s\n", a.String())
	}
}

// Parse will parse the input sources (goargs.DefaultSourceOrder) only once, these values are set for the full program lifetime
func Parse() {
	ParseWith(os.Args[1:], os.Environ())
}

func parseFile() {
	//
}

func lookupEnv(envs []string, key string) (string, bool) {
	for _, env := range envs {
		if strings.HasPrefix(env, key + "=") {
			return strings.SplitN(env, "=", 2)[1], true
		}
	}
	return "", false
}

func parseEnv(envs []string) {
	for _, argDefinition := range DefaultArgs {
		value, found := lookupEnv(envs, toEnvName(argDefinition.longName))
		if found {
			*(argDefinition.refS) = value
			*(argDefinition.source) = SourceEnv
		}
	}
}

func toEnvName(longName string) string {
	return fmt.Sprintf("%s%s", EnvPrefix, strings.ToUpper(longName))
}

func parseArgs(args []string) {
	for _, argDefinition := range DefaultArgs {
		logf("searching for arg: %s (%s)\n", argDefinition.longName, argDefinition.shortName)

		found := false

		inner:
		for _, arg := range args {
			if found {
				*(argDefinition.refS) = arg
				logf("  -> %s\n", arg)
				break inner
			}
			if 
				arg == "--" + argDefinition.longName ||
				(len(argDefinition.shortName) > 0 && arg == "-" + argDefinition.shortName) {

				*(argDefinition.source) = SourceArgs
				found = true
				continue
			}
		}
	}
}

// Watch will parse the input sources (goargs.DefaultSourceOrder) like file and env regularly to reflect the new values
func Watch() {

}

func ResetValues() {
	values = make(map[string]interface{}, 0)
}

func String() string {
	sb := strings.Builder{}
	for _, a := range DefaultArgs {
		fmt.Fprintf(&sb, " --%s ", a.longName)
		fmt.Fprintf(&sb, "%s", *a.refS)
	}
	if sb.Len() > 0 {
		return sb.String()[1:]
	}
	return ""
}

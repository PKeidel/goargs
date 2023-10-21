package goargs

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// [x] Support strings
// [x] Support ints
// [ ] Support flags (bool)
// [x] Support os.Args parsing
// [x] Support env parsing
// [ ] Support config file parsing
// [ ] Support multi flag to allow arrays like `--host host01 host02`
// [ ] Detect missing args
// [ ] Detect undefined args
// [ ] Detect wrong types
// [ ] Create go generator for Arg struct
// [x] Support require flag

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
	DefaultArgs []*Arg
	ArgsByLongname = make(map[string]*Arg, 0)
	DefaultSourceOrder = []SourceType{SourceFile, SourceEnv, SourceArgs}
	EnvPrefix = "GOARG_"

	sourceNotFound SourceType = SourceNotFound
)

type Arg struct {
	argtype ArgType
	longName, shortName, envName string
	Source *SourceType
	refS *string
	refI *int
	Required bool
}

func (a *Arg) String() string {
	switch a.argtype {
	case TypeInt:
		return fmt.Sprintf(
			"Arg[type: %s, longName: %s, shortName: %s, source: %s, ref: %d]",
			a.argtype,
			a.longName,
			a.shortName,
			a.Source,
			*a.refI,
		)
	case TypeString:
		return fmt.Sprintf(
			"Arg[type: %s, longName: %s, shortName: %s, source: %s, ref: %s]",
			a.argtype,
			a.longName,
			a.shortName,
			a.Source,
			*a.refS,
		)
	}
	return fmt.Sprintf("Arg[type: unknown] => %s", a.argtype.String())
}

func Bind(obj interface{}) {
	val := reflect.ValueOf(obj)

	if val.Kind() != reflect.Ptr {
		panic("Input is not a pointer")
	}

	elem := val.Elem()

	if elem.Kind() != reflect.Struct {
		panic("Input is not a struct")
	}

	typ := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		addr := elem.Field(i).Addr().Interface()

		var (
			arg *Arg
			longName = strings.ToLower(typ.Field(i).Name)
		)

		switch ptr := addr.(type) {
		case *string:
			arg = WithStringS(longName, "", ptr)
		case *int:
			arg = WithIntS(longName, "", ptr)
		default:
			logf("unknown type: %v\n", ptr)
		}

		tags := strings.Split(typ.Field(i).Tag.Get("goargs"), ",")

		for _, kvPair := range tags {
			tagInfo := strings.SplitN(kvPair, ":", 2)
			if len(tagInfo[0]) == 0 {
				continue
			}

			logf("tagInfo: %#v\n", tagInfo)

			if tagInfo[0] == "required" {
				arg.Required = true
			}

			if tagInfo[0] == "long" {
				arg.longName = tagInfo[1]
			}
			if tagInfo[0] == "short" {
				arg.shortName = tagInfo[1]
			}
		}
	}
}

func WithIntS(longName, shortName string, ref *int) *Arg {
	logf("WithIntS(%s, %s, %#v)", longName, shortName, ref)
	arg := &Arg{
		argtype: TypeInt,
		longName: longName,
		shortName: shortName,
		Source: &sourceNotFound,
		refI: ref,
		Required: false,
	}
	DefaultArgs = append(DefaultArgs, arg)
	ArgsByLongname[longName] = arg
	return arg
}

func WithStringS(longName, shortName string, ref *string) *Arg {
	logf("WithStringS(%s, %s, %#v)", longName, shortName, ref)
	arg := &Arg{
		argtype: TypeString,
		longName: longName,
		shortName: shortName,
		Source: &sourceNotFound,
		refS: ref,
		Required: false,
	}
	DefaultArgs = append(DefaultArgs, arg)
	ArgsByLongname[longName] = arg
	return arg
}

func WithString(longName string, ref *string) *Arg {
	return WithStringS(longName, "", ref)
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
}

// Parse will parse the input sources (goargs.DefaultSourceOrder) only once, these values are set for the full program lifetime
func Parse() {
	ParseWith(os.Args[1:], os.Environ())
}

func parseFile() {
	log("parseFile()")
	// TODO inotify
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
	log("parseEnv()")
	for _, argDefinition := range DefaultArgs {
		if len(argDefinition.envName) == 0 {
			argDefinition.envName = toEnvName(argDefinition.longName)
		}
		value, found := lookupEnv(envs, argDefinition.envName)
		if found {
			*(argDefinition.refS) = value
			*(argDefinition.Source) = SourceEnv
		}
	}
}

func toEnvName(longName string) string {
	return fmt.Sprintf("%s%s", EnvPrefix, strings.ToUpper(longName))
}

func parseArgs(args []string) {
	log("parseArgs()")
	for _, argDefinition := range DefaultArgs {
		logf("searching for arg: %s (%s)\n", argDefinition.longName, argDefinition.shortName)

		readNextAsValue := false

		inner:
		for _, arg := range args {
			if readNextAsValue {
				readNextAsValue = false
				switch argDefinition.argtype {
				case TypeInt:
					i, err := strconv.Atoi(arg)
					if err != nil {
						panic(err)
					}
					*(argDefinition.refI) = i
					logf("  -> %d\n", i)
				case TypeString:
					*(argDefinition.refS) = arg
					logf("  -> %s\n", arg)
				}
				break inner
			}
			if 
				arg == "--" + argDefinition.longName ||
				(len(argDefinition.shortName) > 0 && arg == "-" + argDefinition.shortName) {

				*(argDefinition.Source) = SourceArgs
				readNextAsValue = true
				continue
			}
		}

		if argDefinition.Required && !readNextAsValue {
			panic("arg is required but was not provided: " + argDefinition.longName)
		}
	}
}

// Watch will parse the input sources (goargs.DefaultSourceOrder) like file and env regularly to reflect the new values
// func Watch() {

// }

func String() string {
	sb := strings.Builder{}
	for _, a := range DefaultArgs {
		fmt.Fprintf(&sb, " --%s ", a.longName)
		switch a.argtype {
		case TypeString:
			fmt.Fprintf(&sb, "%s", *a.refS)
		case TypeInt:
			fmt.Fprintf(&sb, "%d", *a.refI)
		default:
			fmt.Fprint(&sb, " !! unknown type !!")
		}
	}
	if sb.Len() > 0 {
		return sb.String()[1:]
	}
	return ""
}

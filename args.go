package goargs

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

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
	DefaultArgs        []*Arg
	ArgsByLongname     = make(map[string]*Arg, 0)
	DefaultSourceOrder = []SourceType{SourceFile, SourceEnv, SourceArgs}
	EnvPrefix          = "GOARG_"
)

type Arg struct {
	Argtype                      ArgType
	LongName, ShortName, EnvName string
	Source                       SourceType
	refS                         *string
	refI                         *int
	refB                         *bool
	Required                     bool
}

func (a *Arg) HasVal() bool {
	return a.refS != nil || a.refI != nil || a.refB != nil
}

func (a *Arg) Val() interface{} {
	switch a.Argtype {
	case TypeInt:
		return *a.refI
	case TypeString:
		return *a.refS
	case TypeBool:
		return *a.refB
	}
	return "?"
}

func (a *Arg) String() string {
	return fmt.Sprintf(
		"Arg[type: %s, longName: %s, shortName: %s, source: %s, ref: %v]",
		a.Argtype,
		a.LongName,
		a.ShortName,
		a.Source,
		a.Val(),
	)
}

func reflectInfos(obj interface{}) (elem reflect.Value, typ reflect.Type) {
	val := reflect.ValueOf(obj)

	if val.Kind() != reflect.Ptr {
		panic("Input is not a pointer")
	}

	elem = val.Elem()

	if elem.Kind() != reflect.Struct {
		panic("Input is not a struct")
	}

	typ = elem.Type()
	return
}

func Bind(obj interface{}) {
	elem, typ := reflectInfos(obj)

	for i := 0; i < elem.NumField(); i++ {
		addr := elem.Field(i).Addr().Interface()

		var (
			arg      *Arg
			longName = strcase.ToKebab(typ.Field(i).Name)
		)

		switch ptr := addr.(type) {
		case *string:
			arg = WithStringS(longName, "", ptr)
		case *int:
			arg = WithIntS(longName, "", ptr)
		case *bool:
			arg = WithBoolS(longName, "", ptr)
		default:
			logf("unknown type: %v\n", ptr)
		}

		tags := strings.Split(typ.Field(i).Tag.Get("goargs"), ",")

		applyTags(arg, tags)

		arg.EnvName = toEnvName(arg.LongName)
	}
}

func applyTags(arg *Arg, tags []string) {
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
			arg.LongName = tagInfo[1]
		}
		if tagInfo[0] == "short" {
			arg.ShortName = tagInfo[1]
		}
	}
}

func NewDefaultArg(argtype ArgType, longName, shortName string) *Arg {
	return &Arg{
		Argtype:   argtype,
		LongName:  longName,
		ShortName: shortName,
		Source:    SourceNotFound,
		Required:  false,
	}
}

func WithBoolS(longName, shortName string, ref *bool) *Arg {
	logf("WithBoolS(%s, %s, %#v)", longName, shortName, ref)
	arg := NewDefaultArg(TypeBool, longName, shortName)
	arg.refB = ref
	DefaultArgs = append(DefaultArgs, arg)
	ArgsByLongname[longName] = arg
	return arg
}

func WithIntS(longName, shortName string, ref *int) *Arg {
	logf("WithIntS(%s, %s, %#v)", longName, shortName, ref)
	arg := NewDefaultArg(TypeInt, longName, shortName)
	arg.refI = ref
	DefaultArgs = append(DefaultArgs, arg)
	ArgsByLongname[longName] = arg
	return arg
}

func WithStringS(longName, shortName string, ref *string) *Arg {
	logf("WithStringS(%s, %s, %#v)", longName, shortName, ref)
	arg := NewDefaultArg(TypeString, longName, shortName)
	arg.refS = ref
	DefaultArgs = append(DefaultArgs, arg)
	ArgsByLongname[longName] = arg
	return arg
}

func WithString(longName string, ref *string) *Arg {
	return WithStringS(longName, "", ref)
}

func ParseWith(args, envs []string) {
	for _, source := range DefaultSourceOrder {
		switch source {
		case SourceFile:
			parseFile()
		case SourceEnv:
			parseEnv(envs)
		case SourceArgs:
			parseArgs(args)
		}
	}
	checkRequired()
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
		if strings.HasPrefix(env, key+"=") {
			return strings.SplitN(env, "=", 2)[1], true
		}
	}
	return "", false
}

func parseEnv(envs []string) {
	log("parseEnv()")
	for _, argDefinition := range DefaultArgs {
		value, found := lookupEnv(envs, argDefinition.EnvName)
		if found {
			*(argDefinition.refS) = value
			argDefinition.Source = SourceEnv
			fmt.Printf("Set source for %s to %s\n", argDefinition.LongName, SourceEnv.String())
		}
	}
}

func toEnvName(longName string) string {
	return fmt.Sprintf("%s%s", EnvPrefix, strings.ToUpper(longName))
}

func parseArgs(args []string) {
	log("parseArgs()")
	for _, argDefinition := range DefaultArgs {
		logf("searching for arg: %s (%s)\n", argDefinition.LongName, argDefinition.ShortName)

		readNextAsValue := false

	inner:
		for _, arg := range args {
			if readNextAsValue {
				readNextAsValue = false
				switch argDefinition.Argtype {
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
			if arg == "--"+argDefinition.LongName ||
				(len(argDefinition.ShortName) > 0 && arg == "-"+argDefinition.ShortName) {

				argDefinition.Source = SourceArgs

				if argDefinition.Argtype == TypeBool {
					*(argDefinition.refB) = true
					continue
				}

				readNextAsValue = true

				continue
			}
		}
	}
}

func checkRequired() {
	for _, argDefinition := range DefaultArgs {
		if argDefinition.Required && argDefinition.Source == SourceNotFound {
			panic("arg is required but was not provided: " + argDefinition.LongName)
		}
	}
}

// Watch will parse the input sources (goargs.DefaultSourceOrder) like file and env regularly to reflect the new values
// func Watch() {

// }

func String() string {
	sb := strings.Builder{}
	for _, a := range DefaultArgs {
		if a.HasVal() {
			if a.Argtype == TypeBool && a.Val().(bool) == true {
				fmt.Fprintf(&sb, " --%s", a.LongName)
			} else {
				fmt.Fprintf(&sb, " --%s %v", a.LongName, a.Val())
			}
		}
	}
	if sb.Len() > 0 {
		return sb.String()[1:]
	}
	return ""
}

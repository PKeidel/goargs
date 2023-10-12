package goargs

import (
	"testing"
)

func testsetup() {
	log("setup()")
	DefaultArgs = []Arg{}
	DefaultSourceOrder = []SourceType{SourceFile, SourceEnv, SourceArgs}
	ResetValues()
}

func TestSingeLongCmdArgument(t *testing.T) {
	testsetup()

	type Args struct {
		host string
	}
	var (
		args Args
		expectedArgs = "--host localhost"
	)

	EnvPrefix = "TEST_"

	WithStringS("host", "h", &args.host)

    ParseWith(
		[]string{"test", "--host", "localhost"},
		[]string{"TEST_USER=PKeidel"},
	)

	if args.host != "localhost" {
		t.Logf("Expected `args.host` to have value 'localhost', but found: '%s'", args.host)
		t.Fail()
	}

    if String() != expectedArgs {
		t.Logf("Expected `goargs.String()` to return value '%s', but found: '%s'", expectedArgs, String())
		t.Fail()
	}
}

func TestSingeShortCmdArgument(t *testing.T) {
	testsetup()
	
	type Args struct {
		host string
	}
	var (
		args Args
		expectedArgs = "--host localhost"
	)

	WithStringS("host", "h", &args.host)

    ParseWith(
		[]string{"test", "-h", "localhost"},
		[]string{"TEST_USER=PKeidel"},
	)

	if args.host != "localhost" {
		t.Logf("Expected `args.host` to have value 'localhost', but found: '%s'", args.host)
		t.Fail()
	}

    if String() != expectedArgs {
		t.Logf("Expected `goargs.String()` to return value '%s', but found: '%s'", expectedArgs, String())
		t.Fail()
	}
}

func TestSourceSortOrderEnvArgs(t *testing.T) {
	testsetup()
	DefaultSourceOrder = []SourceType{SourceEnv, SourceArgs}

	type Args struct {
		host string
	}
	var args Args

	EnvPrefix = "TEST_"

	WithStringS("host", "h", &args.host)

    ParseWith(
		[]string{"-h", "hostfromargs"},
		[]string{"TEST_HOST=hostfromenv"},
	)

	if args.host != "hostfromargs" {
		t.Logf("Expected `args.host` to have value 'hostfromargs', but found: '%s'", args.host)
		t.Fail()
	}
}

func TestSourceSortOrderArgsEnv(t *testing.T) {
	testsetup()
	DefaultSourceOrder = []SourceType{SourceArgs, SourceEnv}

	type Args struct {
		host string
	}
	var args Args

	EnvPrefix = "TEST_"

	WithStringS("host", "h", &args.host)

    ParseWith(
		[]string{"-h", "hostfromargs"},
		[]string{"TEST_HOST=hostfromenv"},
	)

	if args.host != "hostfromenv" {
		t.Logf("Expected `args.host` to have value 'hostfromenv', but found: '%s'", args.host)
		t.Fail()
	}
}

func TestSingeEnv(t *testing.T) {
	testsetup()

	type Args struct {
		host string
	}
	var (
		args Args
		expectedArgs = "--host localhost"
	)

	EnvPrefix = "TEST_"

	WithStringS("host", "h", &args.host)

    ParseWith(
		[]string{},
		[]string{"TEST_HOST=localhost"},
	)

	if args.host != "localhost" {
		t.Logf("Expected `args.host` to have value 'localhost', but found: '%s'", args.host)
		t.Fail()
	}

    if String() != expectedArgs {
		t.Logf("Expected `goargs.String()` to return value '%s', but found: '%s'", expectedArgs, String())
		t.Fail()
	}
}

func TestALotShortCmdArguments(t *testing.T) {
	testsetup()
	
	type Args struct {
		host, port, user, password string
	}
	var (
		args Args
		expectedArgs = "--host localhost --port 1234"
	)

	EnvPrefix = "TEST_"

	WithStringS("host", "h", &args.host)
	WithStringS("port", "p", &args.port)

    ParseWith(
		[]string{"test", "-h", "localhost", "--port", "1234"},
		[]string{"TEST_USER=PKeidel"},
	)

	if args.host != "localhost" {
		t.Logf("Expected `args.host` to have value '%s', but found: '%s'", "localhost", args.host)
		t.Fail()
	}

    if String() != expectedArgs {
		t.Logf("Expected `goargs.String()` to return value '%s', but found: '%s'", expectedArgs, String())
		t.Fail()
	}
}

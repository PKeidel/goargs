package goargs

import (
	"testing"
)

type ArgsHost struct {
	Host string `goargs:"long:host,short:h"`
}
type ArgsUser struct {
	User string `goargs:"long:user"`
}
type ArgsPort struct {
	Port int `goargs:"short:p"`
}

func testsetup() {
	log("setup()")
	EnvPrefix = "GOARG_"
	DefaultArgs = []*Arg{}
	DefaultSourceOrder = []SourceType{SourceFile, SourceEnv, SourceArgs}
}

func checkString(t *testing.T, expected string) {
	if String() != expected {
		t.Errorf("Expected `goargs.String()` to return value '%s', but found: '%s'", expected, String())
	}
}

func BenchmarkTest(b *testing.B) {
	testsetup()

	var (
		argsHost ArgsHost
	)
	Bind(&argsHost)

	for i := 0; i < b.N; i++ {
		ParseWith(
			[]string{"test", "--host", "localhost"},
			[]string{"TEST_USER=PKeidel"},
		)
	}
}

func TestSingeLongStringCmdArgument(t *testing.T) {
	testsetup()

	var (
		expectedArgs = "--host localhost"
		argsHost     ArgsHost
	)

	Bind(&argsHost)

	ParseWith(
		[]string{"test", "--host", "localhost"},
		[]string{"TEST_USER=PKeidel"},
	)

	checkString(t, expectedArgs)
}

func TestSingeLongIntCmdArgument(t *testing.T) {
	testsetup()

	var expectedArgs = "--port 2345"

	var argsPort ArgsPort

	Bind(&argsPort)

	ParseWith(
		[]string{"test", "--port", "2345"},
		[]string{},
	)

	checkString(t, expectedArgs)
}

func TestSingeShortCmdArgument(t *testing.T) {
	testsetup()

	var (
		expectedArgs = "--host localhost"
		argsHost     ArgsHost
	)

	Bind(&argsHost)

	ParseWith(
		[]string{"test", "-h", "localhost"},
		[]string{"TEST_USER=PKeidel"},
	)

	checkString(t, expectedArgs)
}

func TestSourceSortOrderEnvArgs(t *testing.T) {
	testsetup()
	DefaultSourceOrder = []SourceType{SourceEnv, SourceArgs}
	EnvPrefix = "TEST_"

	var (
		expectedArgs = "--host hostfromargs"
		argsHost     ArgsHost
	)

	Bind(&argsHost)

	ParseWith(
		[]string{"-h", "hostfromargs"},
		[]string{"TEST_HOST=hostfromenv"},
	)

	checkString(t, expectedArgs)
}

func TestSourceSortOrderArgsEnv(t *testing.T) {
	testsetup()
	DefaultSourceOrder = []SourceType{SourceArgs, SourceEnv}
	EnvPrefix = "TEST_"

	var (
		expectedArgs = "--host hostfromenv"
		argsHost     ArgsHost
	)

	Bind(&argsHost)

	ParseWith(
		[]string{"-h", "hostfromargs"},
		[]string{"TEST_HOST=hostfromenv"},
	)

	checkString(t, expectedArgs)
}

func TestSingeEnv(t *testing.T) {
	testsetup()
	EnvPrefix = "TEST_"

	var (
		expectedArgs = "--host hostfromenv"
		argsHost     ArgsHost
	)

	Bind(&argsHost)

	ParseWith(
		[]string{},
		[]string{"TEST_HOST=hostfromenv"},
	)

	checkString(t, expectedArgs)
}

func TestALotShortCmdArguments(t *testing.T) {
	testsetup()
	EnvPrefix = "TEST_"

	var (
		expectedArgs = "--host localhost --port 1234 --user PKeidel"
		argsHostPort struct {
			Host string `goargs:"long:host,short:h"`
			Port int    `goargs:"short:p"`
			User string `goargs:"long:user"`
		}
	)

	Bind(&argsHostPort)

	ParseWith(
		[]string{"test", "-h", "localhost", "-p", "1234"},
		[]string{"TEST_USER=PKeidel"},
	)

	checkString(t, expectedArgs)
}

func TestRequiredArguments(t *testing.T) {
	defer func() {
		// we want a panic()! So if there is no error, the test should fail
		if r := recover(); r == nil {
			t.Fatal("Test must panic()! Because a required arg was not provided")
		}
	}()

	testsetup()

	var (
		argsHost struct {
			Host string `goargs:"long:host,short:h,required"`
		}
	)

	Bind(&argsHost)

	ParseWith(
		[]string{"test", "-x", "localhost"},
		[]string{},
	)
}

func TestRequiredEnv1(t *testing.T) {
	testsetup()

	var (
		expectedArgs = "--host localhost"
		argsHost     struct {
			Host string `goargs:"required"`
		}
	)

	EnvPrefix = ""
	Bind(&argsHost)

	ParseWith(
		[]string{},
		[]string{"HOST=localhost"},
	)

	checkString(t, expectedArgs)
}

func TestRequiredEnv2(t *testing.T) {
	testsetup()

	var (
		expectedArgs = "--host localhost"
		argsHost     struct {
			Host string `goargs:"required"`
		}
	)

	Bind(&argsHost)

	ParseWith(
		[]string{},
		[]string{"GOARG_HOST=localhost"},
	)

	checkString(t, expectedArgs)
}

func TestSingleVar01(t *testing.T) {
	testsetup()

	var (
		expectedArgs = "--host localhost"
		host         string
	)
	WithStringS("host", "h", &host)

	ParseWith(
		[]string{"test", "-h", "localhost"},
		[]string{"TEST_USER=PKeidel"},
	)

	checkString(t, expectedArgs)
}

func TestBool1(t *testing.T) {
	testsetup()

	var (
		expectedArgs = "--host localhost --debug"
		argHostDebug struct {
			Host  string `goargs:"long:host,short:h"`
			Debug bool
		}
	)

	Bind(&argHostDebug)

	ParseWith(
		[]string{"test", "-h", "localhost", "--debug"},
		[]string{"TEST_USER=PKeidel"},
	)

	checkString(t, expectedArgs)
}

func TestExternalChange1(t *testing.T) {
	testsetup()

	var (
		expectedArgs = "--host localhost"
		argHostDebug struct {
			Host string `goargs:"long:host"`
		}
	)

	Bind(&argHostDebug)

	ParseWith(
		[]string{},
		[]string{},
	)

	// now we modify the value directly
	argHostDebug.Host = "localhost"

	// the new value should be included in the generated string
	checkString(t, expectedArgs)
}

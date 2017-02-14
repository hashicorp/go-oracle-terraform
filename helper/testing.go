package helper

import (
	"fmt"
	"os"
)

const TestEnvVar = "ORACLE_ACC"

// Test suite helpers

type TestCase struct {
	// Fields to test stuff with
}

func Test(t TestT, c TestCase) {
	if os.Getenv(TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' is set", TestEnvVar))
		return
	}
}

type TestT interface {
	Error(args ...interface{})
	Fatal(args ...interface{})
	Skip(args ...interface{})
}

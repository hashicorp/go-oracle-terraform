// Copyright IBM Corp. 2017, 2026

package httpmock

import (
	"os"
)

var envVarName = "GONOMOCKS"

func Disabled() bool {
	return os.Getenv(envVarName) != ""
}

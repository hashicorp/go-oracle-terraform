// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package httpmock

import (
	"os"
)

var envVarName = "GONOMOCKS"

func Disabled() bool {
	return os.Getenv(envVarName) != ""
}

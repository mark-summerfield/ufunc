// Copyright © 2024 Mark Summerfield. All rights reserved.

package ufunc

import (
    "fmt"
    _ "embed"
    )

//go:embed Version.dat
var Version string

func Hello() string {
    return fmt.Sprintf("Hello ufunc v%s", Version)
}

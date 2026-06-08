package version

import (
	"fmt"
	"runtime"
)

// These variables are populated at build time using -ldflags
var (
	Version   = "dev"
	BuildTime = "unknown"
)

// Info returns a formatted string containing full build details.
func Info() string {
	return fmt.Sprintf("v%s (%s) built on %s using %s", Version, BuildTime, runtime.Version())
}

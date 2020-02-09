package version

import (
	"fmt"
)

var (
	version = "0.0.0"
	commit  = "HEAD"
)

func Version() string {
	return version
}

func Commit() string {
	return commit
}

func Application(application string) string {
	return fmt.Sprintf("%s (%s-%s)", application, version, commit)
}

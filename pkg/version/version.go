package version

import (
	"fmt"
)

var (
	version = "0.0.0"
	commit  = "HEAD"
	date    = "1970-01-01T00:00:00Z00:00"
)

func Version() string {
	return version
}

func Commit() string {
	return commit
}

func Date() string {
	return date
}

func Application(application string) string {
	return fmt.Sprintf("%s (%s-%s) [%s]", application, version, commit, date)
}

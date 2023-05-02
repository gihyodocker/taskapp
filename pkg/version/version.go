package version

import (
	"fmt"
)

var (
	version = "unspecified"
)

type Info struct {
	Version string
}

func Get() Info {
	return Info{
		Version: version,
	}
}

func (i Info) String() string {
	return fmt.Sprintf("taskapp version: %s", i.Version)
}

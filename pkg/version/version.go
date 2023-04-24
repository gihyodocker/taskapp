package version

import (
	"fmt"
	"runtime/debug"
)

const (
	unspecified = "unspecified"
)

var (
	osArch    = unspecified
	gitCommit = unspecified
	version   = unspecified
)

type Info struct {
	OSArch    string
	GitCommit string
	Version   string
}

func Get() Info {
	return Info{
		OSArch:    osArch,
		GitCommit: gitCommit,
		Version:   version,
	}
}

func (i Info) String() string {
	if i.Version == unspecified {
		info, _ := debug.ReadBuildInfo()
		return fmt.Sprintf("todoapp Version: %s", info.Main.Version)
	}

	return fmt.Sprintf(
		"todoapp %s Version: %s, GitCommit: %s",
		i.OSArch,
		i.Version,
		i.GitCommit,
	)
}

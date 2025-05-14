package version

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

type Info struct {
	Version   string
	BuildTime string
	GitCommit string
}

func Get() Info {
	return Info{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
	}
}

package version

var (
	// Version is the current version of the application
	Version = "dev"
	// BuildTime is the time when the application was built
	BuildTime = "unknown"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
)

// Info contains version information
type Info struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
}

// Get returns the version information
func Get() Info {
	return Info{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
	}
}

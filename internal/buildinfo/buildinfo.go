package buildinfo

var (
	Version = "dev"
	Commit  = ""
	Date    = ""
)

func Summary() string {
	version := Version
	if version == "" {
		version = "dev"
	}
	commit := Commit
	if commit == "" {
		commit = "unknown"
	}
	date := Date
	if date == "" {
		date = "unknown"
	}
	return "ooxmlx " + version + " (" + commit + ", " + date + ")"
}

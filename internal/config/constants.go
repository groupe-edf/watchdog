package config

// Banner watchdog CLI banner
const Banner = `
┬ ┬┌─┐┌┬┐┌─┐┬ ┬┌┬┐┌─┐┌─┐
│││├─┤ │ │  ├─┤ │││ ││ ┬
└┴┘┴ ┴ ┴ └─┘┴ ┴─┴┘└─┘└─┘

Watchdog Version v{{ .BuildInfo.Version }} build {{ .BuildInfo.Sha }}
{{ if .Options.Contact }}{{ .Options.Contact }}{{ end }}
{{ if .Options.DocsLink }}Documentation: {{ .Options.DocsLink }}{{ end }}
`

var (
	// LogPath default log file location
	LogPath = "/var/log/watchdog/watchdog.log"
	// ErrorMessagePrefix To have custom error messages appear in GitLab’s UI when the commit is declined or an error occurs during the Git hook
	ErrorMessagePrefix = "GL-HOOK-ERR: "
	// BallotX failure icon
	BallotX = "\u2717"
	// CheckMark success mark
	CheckMark = "\u2713"
)

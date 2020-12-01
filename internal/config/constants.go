package config

// Banner watchdog CLI banner
const Banner = `
┬ ┬┌─┐┌┬┐┌─┐┬ ┬┌┬┐┌─┐┌─┐
│││├─┤ │ │  ├─┤ │││ ││ ┬
└┴┘┴ ┴ ┴ └─┘┴ ┴─┴┘└─┘└─┘

Watchdog Version v{{ .BuildInfo.Version }}+{{ printf "%.8s" .BuildInfo.Commit }}
{{ if .Options.Contact -}}
Contact: {{ .Options.Contact }}
{{ end -}}
{{ if .Options.DocsLink -}}
Documentation: {{ .Options.DocsLink }}
{{ end -}}
`

var (
	// BallotX failure icon
	BallotX = "\u2717"
	// CheckMark success mark
	CheckMark = "\u2713"
)

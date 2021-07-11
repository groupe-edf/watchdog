package output

type Format string

const (
	CSV  Format = "csv"
	JSON        = "json"
	Text        = "text"
)

type Options struct {
	// Used to print error messages with high severity in some Git UI.
	// Gitlab for instance use `GL-HOOK-ERR`. See https://docs.gitlab.com/ee/administration/server_hooks.html
	ErrorMessagePrefix string
	Format             Format
	Template           string
}

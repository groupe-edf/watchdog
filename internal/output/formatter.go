package output

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/util"
)

var text = `
{{- if .Issues -}}
{{ range .Issues -}}
{{ $.LinePrefix }}severity={{ .Severity }} handler={{ .Handler }} condition={{ .Condition }} commit={{ printf "%.8s" .Commit }} message="{{ .Message }}"
{{ end -}}
{{ end -}}
`

// ReportData report data
type ReportData struct {
	Issues     []issue.Issue
	LinePrefix string
}

// NewReport return analysis report
func NewReport(writer io.Writer, format string, set *util.Set) (err error) {
	switch format {
	case "json":
		raw, err := json.MarshalIndent(set.List(), "", "\t")
		if err != nil {
			return err
		}
		_, err = writer.Write(raw)
		return err
	case "text":
		functionsMap := template.FuncMap{
			"ToUpper": strings.ToUpper,
		}
		t := template.Must(template.New("watchdog").Funcs(functionsMap).Parse(text))
		return t.Execute(writer, &ReportData{
			Issues:     set.List(),
			LinePrefix: config.ErrorMessagePrefix,
		})
	}
	return errors.New("Unsupported output format")
}

// Report output issues report
func Report(path string, format string, set *util.Set) (err error) {
	if set.Len() > 0 {
		if path != "" {
			file, err := os.Create(path)
			if err != nil {
				return err
			}
			defer file.Close()
			return NewReport(file, format, set)
		}
		os.Stdout.Write([]byte("-----BEGIN REJECTION MESSAGES-----\n"))
		err = NewReport(os.Stdout, format, set)
		os.Stdout.Write([]byte("\n-----BEGIN REJECTION MESSAGES-----"))
		return err
	}
	return nil
}

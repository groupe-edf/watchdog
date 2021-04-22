package output

import (
	"encoding/json"
	"errors"
	"fmt"
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
{{- if eq .Severity 2 -}}{{ $.ErrorMessagePrefix }}{{ end -}}severity={{ .Severity }} handler={{ .Handler }} condition={{ .Condition }} commit={{ printf "%.8s" .Commit }} message="{{ .Message }}"
{{ end -}}
{{ end -}}
`

// ReportData report data
type ReportData struct {
	Issues             []issue.Issue
	ErrorMessagePrefix string
}

// NewReport return analysis report
func NewReport(writer io.Writer, options *config.Options, set *util.Set) (err error) {
	switch options.OutputFormat {
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
			Issues:             set.List(),
			ErrorMessagePrefix: options.ErrorMessagePrefix,
		})
	default:
		return errors.New("unsupported output format")
	}
}

// Report output issues report
func Report(options *config.Options, set *util.Set) (err error) {
	if set.Len() > 0 {
		if options.Output != "" {
			file, err := os.Create(options.Output)
			if err != nil {
				return err
			}
			defer file.Close()
			err = NewReport(file, options, set)
			if err != nil {
				return err
			}
			fmt.Printf("Report file generated in %s", util.Colorize(util.Green, options.Output))
			fmt.Println()
			return nil
		}
		os.Stdout.Write([]byte("-----BEGIN REJECTION MESSAGES-----\n"))
		err = NewReport(os.Stdout, options, set)
		os.Stdout.Write([]byte("\n-----BEGIN REJECTION MESSAGES-----"))
		return err
	}
	return nil
}

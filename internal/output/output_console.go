package output

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/version"
	"github.com/groupe-edf/watchdog/pkg/container"
)

type Console struct {
	Channel chan models.AnalysisResult
	Options Options
	writer  io.Writer
}

const outputTemplate = `
{{- if .Issues -}}
{{ range .Issues -}}
· {{ if eq .Severity 2 -}}{{ $.ErrorMessagePrefix }}{{ end -}}severity={{ .Severity }} handler={{ .PolicyType }} condition={{ .ConditionType }} commit={{ printf "%.8s" .Commit.Hash }} message="{{ .Message }}"
{{ end -}}
{{ end -}}
`

func (output *Console) WriteTo() error {
	output.printBanner()
	severity := models.SeverityLow
	for {
		if result, ok := <-output.Channel; ok {
			for _, data := range result.Issues {
				if data.Severity > severity {
					severity = data.Severity
				}
			}
			var content bytes.Buffer
			var err error
			switch output.Options.Format {
			case CSV:
				writer := csv.NewWriter(os.Stdout)
				for _, issue := range result.Issues {
					offender, _ := json.Marshal(issue.Offender)
					writer.Write([]string{
						issue.Severity.String(),
						fmt.Sprintf("%.8s", issue.Commit.Hash),
						string(issue.PolicyType),
						issue.Message,
						string(offender),
					})
				}
				writer.Flush()
			case JSON:
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "\t")
				for _, issue := range result.Issues {
					err := encoder.Encode(issue)
					if err != nil {
						return err
					}
				}
			case Text:
				commitHash := color.Green.Sprint(result.Commit.Hash[:8])
				if len(result.Issues) > 0 {
					commitHash = color.Red.Sprint(result.Commit.Hash[:8])
				}
				content.Write([]byte(fmt.Sprintf("%v · %v · (%v)\n", commitHash, strings.Split(result.Commit.Subject, "\n")[0], result.ElapsedTime)))
				functionsMap := template.FuncMap{
					"upper": strings.ToUpper,
				}
				t := template.Must(template.New("watchdog").Funcs(functionsMap).Parse(output.Options.Template))
				err = t.Execute(&content, &ReportData{
					Issues:             result.Issues,
					ErrorMessagePrefix: output.Options.ErrorMessagePrefix,
				})
				if err != nil {
					return err
				}
			default:
				return errors.New("unsupported output format")
			}
			_, err = output.writer.Write(content.Bytes())
			if err != nil {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

// printBanner print watchdog banner
func (output *Console) printBanner() error {
	options := container.Get(config.ServiceName).(*config.Options)
	t, err := template.New("watchdog").Parse(config.Banner)
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"Options":   options,
		"BuildInfo": version.GetBuildInfo(),
	}
	return t.Execute(output.writer, data)
}

func NewConsole(ch chan models.AnalysisResult) *Console {
	var writer io.Writer = os.Stdout
	options := container.GetContainer().Get(config.ServiceName).(*config.Options)
	if options.Output != "" {
		writer, err := os.Create(options.Output)
		if err != nil {
			panic(err)
		}
		defer writer.Close()
	}
	return &Console{
		Channel: ch,
		Options: Options{
			ErrorMessagePrefix: options.ErrorMessagePrefix,
			Format:             Format(options.OutputFormat),
			Template:           outputTemplate,
		},
		writer: writer,
	}
}

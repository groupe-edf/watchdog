package core

import (
	"context"
	"testing"

	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/core/handlers"
	"github.com/groupe-edf/watchdog/internal/git"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/server/container"
)

func BenchmarkTestAnalyze(b *testing.B) {
	b.StopTimer()
	container.GetContainer().Provide(&logging.ServiceProvider{})
	logger := container.GetContainer().Get(logging.ServiceName).(logging.Interface)
	ctx := context.Background()
	analyzer := &Analyzer{
		context: ctx,
		Logger:  logger,
		Policies: []models.Policy{
			{
				Conditions: []models.Condition{
					{
						Pattern: "",
						Type:    models.ConditionTypePattern,
					},
				},
				Type: models.PolicyTypeCommit,
			},
		},
	}
	analyzer.RegisterHandler(&handlers.CommitHandler{})
	analyzeChan := make(chan models.AnalysisResult)
	client := git.NewGit(&config.Options{})
	_, err := client.Clone(ctx, git.CloneOptions{
		URL: "https://github.com/groupe-edf/watchdog",
	})
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()
	commitIter, _ := client.Commits(context.Background(), git.LogOptions{})
	analyzer.Analyze(commitIter, analyzeChan)
	b.StopTimer()
}

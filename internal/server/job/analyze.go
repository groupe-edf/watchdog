package job

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/google/uuid"
	"github.com/gookit/color"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/core"
	"github.com/groupe-edf/watchdog/internal/core/loaders"
	driver "github.com/groupe-edf/watchdog/internal/git"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/server/broadcast"
	"github.com/groupe-edf/watchdog/internal/server/container"
	"github.com/groupe-edf/watchdog/internal/util"
)

type AnalyzeOptions struct {
	AnalysisID    *uuid.UUID `json:"analysis_id,omitempty"`
	CreatedAt     *time.Time `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by"`
	From          string     `json:"last_commit_hash"`
	IntegrationID int64      `json:"integration_id"`
	RepositoryID  *uuid.UUID `json:"repository_id,omitempty"`
	RepositoryURL string     `json:"repository_url,omitempty"`
}

type ProcessAnalyze struct {
	Context context.Context
	Logger  logging.Interface
	Options *config.Options
	Store   models.Store
}

func (processor *ProcessAnalyze) Handle(job *models.Job) error {
	var options AnalyzeOptions
	if err := json.Unmarshal(job.Args, &options); err != nil {
		return err
	}
	analysis, err := processor.Store.FindAnalysisByID(options.AnalysisID)
	if err != nil {
		return err
	}
	analysis.Start()
	_, err = processor.Store.SaveAnalysis(analysis)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(processor.Context, time.Duration(time.Second*1240))
	defer cancel()
	cloneOptions := driver.CloneOptions{
		URL: options.RepositoryURL,
	}
	if options.IntegrationID != 0 {
		integration, err := processor.Store.FindIntegrationByID(options.IntegrationID)
		if err != nil {
			return err
		}
		token, err := util.Decrypt(integration.APIToken, processor.Options.Server.Security.MasterKey)
		if err != nil {
			return err
		}
		cloneOptions.Authentication = &driver.BasicAuthentication{
			Username: integration.InstanceName,
			Password: token,
		}
	}
	client := container.Get(driver.ServiceName).(driver.Driver)
	_, err = client.Clone(ctx, cloneOptions)
	if err != nil && err != git.NoErrAlreadyUpToDate {
		if err == transport.ErrAuthenticationRequired || err == transport.ErrEmptyRemoteRepository {
			analysis.State = models.FailedState
			analysis.StateMessage = err.Error()
			_, err = processor.Store.SaveAnalysis(analysis)
			return err
		}
		return err
	}
	reference, err := client.Head()
	if err != nil {
		return err
	}
	analysis.LastCommitHash = reference
	from := plumbing.NewHash(options.From)
	var commits models.Iterator[models.Commit]
	if !from.IsZero() {
		commits, err = client.RevList(driver.RevListOptions{
			OldRev: from,
		})
	} else {
		commits, err = client.Commits(ctx, driver.LogOptions{})
	}
	if err != nil {
		return err
	}
	analyzer, err := core.NewAnalyzer(
		processor.Context,
		loaders.NewStoreLoader(processor.Store),
		processor.Logger,
		processor.Options,
		models.Whitelist{
			Files: []string{
				"package-lock.json",
			},
			Paths: []string{
				"node_modules",
			},
		},
	)
	if err != nil {
		return err
	}
	analyzeChan := make(chan models.AnalysisResult)
	broadcast := container.GetContainer().Get(broadcast.ServiceName).(*broadcast.Broadcast)
	severity := models.SeverityLow
	totalIssues := 0
	go analyzer.Analyze(commits, analyzeChan)
	for {
		if result, ok := <-analyzeChan; ok {
			for _, data := range result.Issues {
				totalIssues += len(result.Issues)
				if data.Severity > severity {
					severity = data.Severity
				}
				err := processor.Store.SaveIssue(options.RepositoryID, options.AnalysisID, data)
				if err != nil {
					return err
				}
				if len(data.Leaks) > 0 {
					err = processor.Store.SaveLeaks(options.RepositoryID, options.AnalysisID, data.Leaks)
					if err != nil {
						return err
					}
				}
			}
			commitHash := color.Green.Sprint(result.Commit.Hash[:8])
			if len(result.Issues) > 0 {
				commitHash = color.Red.Sprint(result.Commit.Hash[:8])
			}
			fmt.Printf("|_ %v · %v · (%v)\n", commitHash, strings.Split(result.Commit.Subject, "\n")[0], result.ElapsedTime)
			broadcast.Broadcast(result)
		} else {
			break
		}
	}
	analysis.Complete(severity, totalIssues)
	processor.Store.SaveAnalysis(analysis)
	broadcast.Broadcast(map[string]string{
		"event":        "analysis_finished",
		"container_id": analysis.ID.String(),
	})
	return nil
}

package job

import (
	"context"
	"encoding/json"
	"errors"
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
	"github.com/groupe-edf/watchdog/internal/core/models"
	driver "github.com/groupe-edf/watchdog/internal/git"
	"github.com/groupe-edf/watchdog/internal/server/queue"
	"github.com/groupe-edf/watchdog/internal/util"
	"github.com/groupe-edf/watchdog/pkg/container"
	"github.com/groupe-edf/watchdog/pkg/event"
	"github.com/groupe-edf/watchdog/pkg/logging"
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
	publisher := container.GetContainer().Get(event.ServiceName).(*event.EventBus)
	var options AnalyzeOptions
	if err := json.Unmarshal(job.Args, &options); err != nil {
		return err
	}
	analysis, err := processor.Store.FindAnalysisByID(options.AnalysisID)
	if err != nil {
		return err
	}
	analysis.Start()
	defer func() {
		if err := recover(); err != nil {
			if job.ErrorCount == queue.MaxErroCount {
				analysis.Failed(errors.New(job.LastError))
				publisher.PublishAsync("analysis:failed", analysis)
			}
		}
	}()
	publisher.PublishAsync("analysis:started", analysis)
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
	repository, err := client.Clone(ctx, cloneOptions)
	if err != nil && err != git.NoErrAlreadyUpToDate {
		if err == transport.ErrAuthenticationRequired || err == transport.ErrEmptyRemoteRepository {
			analysis.Failed(err)
			publisher.PublishAsync("analysis:failed", analysis)
			return err
		}
		return err
	}
	head, err := client.Head()
	if err != nil {
		return err
	}
	analysis.LastCommitHash = head
	from := plumbing.NewHash(options.From)
	var commits models.Iterator[models.Commit]
	if !from.IsZero() {
		commits, err = client.RevList(driver.RevListOptions{})
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
		repository,
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
	severity := models.SeverityLow
	totalIssues := 0
	go analyzer.Analyze(repository, commits, analyzeChan)
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
			if result.Commit.Hash != "" {
				commitHash := color.Green.Sprint(result.Commit.Hash[:8])
				if len(result.Issues) > 0 {
					commitHash = color.Red.Sprint(result.Commit.Hash[:8])
				}
				fmt.Printf("|_ %v · %v · (%v)\n", commitHash, strings.Split(result.Commit.Subject, "\n")[0], result.ElapsedTime)
			}
		} else {
			break
		}
	}
	analysis.Complete(severity, totalIssues)
	publisher.PublishAsync("analysis:finished", analysis)
	return nil
}

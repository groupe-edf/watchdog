package core

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/core/handlers"
	"github.com/groupe-edf/watchdog/internal/core/loaders"
	"github.com/groupe-edf/watchdog/internal/git"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/security"
	"github.com/groupe-edf/watchdog/internal/util"
)

// Analyzer git commits analyzer
type Analyzer struct {
	Driver    git.Driver
	Handlers  []handlers.Handler
	Info      *hook.Info
	Issues    *util.Set
	Loader    loaders.Loader
	Logger    logging.Interface
	Options   *config.Options
	Policies  []models.Policy
	Whitelist models.Whitelist
	context   context.Context
}

// Analyze execute analysis
func (analyzer *Analyzer) Analyze(commitIter models.Iterator[models.Commit], analyzeChan chan models.AnalysisResult) error {
	defer close(analyzeChan)
	if len(analyzer.Handlers) == 0 {
		return errors.New("at least one handler must be defined")
	}
	if analyzer.context == nil {
		analyzer.context = context.Background()
	}
	ctx, cancel := context.WithCancel(analyzer.context)
	defer cancel()
	if len(analyzer.Policies) > 0 {
		var wg sync.WaitGroup
		// struct{} is the smallest data type available in Go
		maxWorkers := make(chan struct{}, 4)
		var currentCommit chan *models.Commit = make(chan *models.Commit)
		totalCommit := 0
		go func() {
			defer close(currentCommit)
			commitIter.ForEach(func(commit *models.Commit) error {
				totalCommit++
				currentCommit <- commit
				return nil
			})
		}()
		for commit := range currentCommit {
			wg.Add(1)
			maxWorkers <- struct{}{}
			go func(commit *models.Commit) {
				defer wg.Done()
				defer func() { <-maxWorkers }()
				_ = analyzer.analyze(ctx, commit, analyzeChan)
			}(commit)
		}
		wg.Wait()
		close(maxWorkers)
	} else {
		analyzer.Logger.WithFields(logging.Fields{
			"correlation_id": util.GetRequestID(ctx),
			"user_id":        util.GetUserID(ctx),
		}).Info("no policies were found")
	}
	return nil
}

// Context returns underlying context
func (analyzer *Analyzer) Context() context.Context {
	return analyzer.context
}

// HasErrors check id set has issues with high severity
func (analyzer *Analyzer) HasErrors() bool {
	for _, item := range analyzer.Issues.List() {
		if item.Severity == models.SeverityHigh {
			return true
		}
	}
	return false
}

// RegisterHandler register git hook handler
func (analyzer *Analyzer) RegisterHandler(handler handlers.Handler) {
	analyzer.Logger.WithFields(logging.Fields{
		"correlation_id": util.GetRequestID(analyzer.context),
		"user_id":        util.GetUserID(analyzer.context),
	}).Debugf("registring handler `%v`", reflect.TypeOf(handler))
	handler.SetLogger(analyzer.Logger)
	if analyzer.Info != nil {
		handler.SetInfo(analyzer.Info)
	}
	handler.SetOptions(analyzer.Options)
	handler.SetDriver(analyzer.Driver)
	analyzer.Handlers = append(analyzer.Handlers, handler)
}

// SetInfo set hooks
func (analyzer *Analyzer) SetInfo(info *hook.Info) {
	analyzer.Info = info
}

// SetLogger set logger
func (analyzer *Analyzer) SetLogger(logger logging.Interface) {
	analyzer.Logger = logger
}

// SetPolicies set policies
func (analyzer *Analyzer) SetPolicies(policies []models.Policy) {
	analyzer.Policies = policies
}

// SetDriver set repository
func (analyzer *Analyzer) SetDriver(driver git.Driver) {
	analyzer.Driver = driver
}

func (analyzer *Analyzer) analyze(ctx context.Context, commit *models.Commit, analyzeChan chan models.AnalysisResult) error {
	scanTimeStart := time.Now()
	issues := make([]models.Issue, 0)
	for _, policy := range analyzer.Policies {
		analyzer.Logger.WithFields(logging.Fields{
			"commit":         commit.Hash,
			"correlation_id": util.GetRequestID(ctx),
			"rule":           policy.Type,
			"user_id":        util.GetUserID(ctx),
		}).Debug("processing hook rule")
		for _, handler := range analyzer.Handlers {
			if handler.GetType() != handlers.HandlerTypeCommits {
				continue
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default: // Prevent from blocking.
			}
			issuesSlice, err := handler.Handle(ctx, commit, policy, analyzer.Whitelist)
			if err != nil {
				return err
			}
			issues = append(issues, issuesSlice...)
		}
	}
	elapsed := time.Since(scanTimeStart)
	analyzeChan <- models.AnalysisResult{
		Commit:      *commit,
		ElapsedTime: elapsed,
		Issues:      issues,
	}
	return ctx.Err()
}

func (analyzer *Analyzer) handleRef(ctx context.Context) {
	issues := make([]models.Issue, 0)
	for _, policy := range analyzer.Policies {
		analyzer.Logger.WithFields(logging.Fields{
			"correlation_id": util.GetRequestID(ctx),
			"rule":           policy.Type,
			"user_id":        util.GetUserID(ctx),
		}).Debug("processing hook rule")
		for _, handler := range analyzer.Handlers {
			if handler.GetType() == handlers.HandlerTypeRefs {
				issuesSlice, _ := handler.Handle(ctx, nil, policy, analyzer.Whitelist)
				issues = append(issues, issuesSlice...)
			}
		}
	}
	analyzer.Issues.Add(issues)
}

// NewAnalyzer instantiate new analyzer
func NewAnalyzer(
	ctx context.Context,
	loader loaders.Loader,
	logger logging.Interface,
	options *config.Options,
	whitelist models.Whitelist,
) (*Analyzer, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	policies, err := loader.LoadPolicies(ctx)
	if err != nil {
		return nil, err
	}
	analyzer := &Analyzer{
		Options:   options,
		Issues:    util.NewSet(),
		Logger:    logger,
		Policies:  policies,
		Whitelist: whitelist,
		context:   ctx,
	}
	// Register handlers
	analyzer.RegisterHandler(&handlers.BranchHandler{})
	analyzer.RegisterHandler(&handlers.CommitHandler{})
	analyzer.RegisterHandler(&handlers.FileHandler{})
	analyzer.RegisterHandler(&handlers.JiraHandler{})
	rules, err := loader.LoadRules(ctx)
	if err != nil {
		return nil, err
	}
	analyzer.RegisterHandler(&handlers.SecurityHandler{
		Scanner: security.NewRegexScanner(logger, rules, whitelist),
	})
	analyzer.RegisterHandler(&handlers.TagHandler{})
	return analyzer, nil
}

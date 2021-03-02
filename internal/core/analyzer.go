package core

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gookit/color"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/util"
)

// Analyzer git commits analyzer
type Analyzer struct {
	GitHooks   *hook.GitHooks
	Handlers   []Handler
	Info       *hook.Info
	Issues     *util.Set
	Logger     logging.Interface
	Options    *config.Options
	Repository *git.Repository
}

// Analyze execute analysis
func (analyzer *Analyzer) Analyze(ctx context.Context, commitIter object.CommitIter) error {
	analyzer.Logger.WithFields(logging.Fields{
		"correlation_id": util.GetRequestID(ctx),
		"repository":     analyzer.Options.URI,
		"user_id":        util.GetUserID(ctx),
	}).Info("starting analysis")
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer commitIter.Close()
	maxWorkers := make(chan struct{}, analyzer.Options.Concurrent)
	if len(analyzer.GitHooks.Hooks) > 0 {
		analyzer.Logger.WithFields(logging.Fields{
			"correlation_id": util.GetRequestID(ctx),
			"user_id":        util.GetUserID(ctx),
		}).Debugf("%v handlers found and %v hooks found", len(analyzer.Handlers), len(analyzer.GitHooks.Hooks))
		var wg sync.WaitGroup
		for {
			commit, err := commitIter.Next()
			if err == object.ErrUnsupportedObject {
				continue
			}
			if err == io.EOF {
				break
			}
			wg.Add(1)
			maxWorkers <- struct{}{}
			go func(commit *object.Commit) {
				defer wg.Done()
				defer func() { <-maxWorkers }()
				err := analyzer.analyze(ctx, analyzer.GitHooks, commit)
				if err != nil {
					analyzer.Logger.WithFields(logging.Fields{
						"commit":         commit.Hash.String(),
						"correlation_id": util.GetRequestID(ctx),
						"user_id":        util.GetUserID(ctx),
					}).Error(err)
				}
			}(commit)
		}
		wg.Wait()
		close(maxWorkers)
	} else {
		analyzer.Logger.WithFields(logging.Fields{
			"correlation_id": util.GetRequestID(ctx),
			"user_id":        util.GetUserID(ctx),
		}).Info("there is no hooks in .githooks.yml file")
	}
	return nil
}

// HasErrors check id set has issues with high severity
func (analyzer *Analyzer) HasErrors() bool {
	for _, item := range analyzer.Issues.List() {
		if item.Severity == issue.SeverityHigh {
			return true
		}
	}
	return false
}

// RegisterHandler register git hook handler
func (analyzer *Analyzer) RegisterHandler(ctx context.Context, handler Handler) {
	analyzer.Logger.WithFields(logging.Fields{
		"correlation_id": util.GetRequestID(ctx),
		"user_id":        util.GetUserID(ctx),
	}).Debugf("registring handler `%v`", reflect.TypeOf(handler))
	handler.SetLogger(analyzer.Logger)
	if analyzer.Info != nil {
		handler.SetInfo(analyzer.Info)
	}
	handler.SetOptions(analyzer.Options)
	handler.SetRepository(analyzer.Repository)
	analyzer.Handlers = append(analyzer.Handlers, handler)
}

// SetHooks set hooks
func (analyzer *Analyzer) SetHooks(hooks *hook.GitHooks) {
	if len(analyzer.Options.Handlers) > 0 {
		for handler, rule := range analyzer.Options.Handlers {
			if rule.Disabled {
				continue
			}
			hooks.Hooks[0].Rules = append(hooks.Hooks[0].Rules, &hook.Rule{
				Conditions:  rule.Conditions,
				Description: rule.Description,
				Disabled:    rule.Disabled,
				Type:        hook.HandlerType(handler),
			})
		}
	}
	analyzer.GitHooks = hooks
}

// SetInfo set hooks
func (analyzer *Analyzer) SetInfo(info *hook.Info) {
	analyzer.Info = info
}

// SetLogger set logger
func (analyzer *Analyzer) SetLogger(logger logging.Interface) {
	analyzer.Logger = logger
}

// SetRepository set repository
func (analyzer *Analyzer) SetRepository(repository *git.Repository) {
	analyzer.Repository = repository
}

func (analyzer *Analyzer) analyze(ctx context.Context, gitHooks *hook.GitHooks, commit *object.Commit) error {
	scanTimeStart := time.Now()
	issues := make([]issue.Issue, 0)
	for _, hook := range gitHooks.Hooks {
		for _, rule := range hook.Rules {
			analyzer.Logger.WithFields(logging.Fields{
				"commit":         commit.Hash.String(),
				"correlation_id": util.GetRequestID(ctx),
				"rule":           rule.Type,
				"user_id":        util.GetUserID(ctx),
			}).Debug("processing hook rule")
			for _, handler := range analyzer.Handlers {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					// Prevent from blocking.
				}
				issuesSlice, _ := handler.Handle(ctx, commit, rule)
				issues = append(issues, issuesSlice...)
			}
		}
	}
	analyzer.Issues.Add(issues)
	commitHash := color.Green.Sprint(commit.Hash.String()[:8])
	if len(issues) > 0 {
		commitHash = color.Red.Sprint(commit.Hash.String()[:8])
	}
	elapsed := time.Since(scanTimeStart)
	fmt.Printf("|_ %v · %v · (%v)\n", commitHash, strings.Split(commit.Message, "\n")[0], elapsed)
	return ctx.Err()
}

// NewAnalyzer instantiate new analyzer
func NewAnalyzer(hooks *hook.GitHooks, options *config.Options) (*Analyzer, error) {
	analyzer := &Analyzer{
		GitHooks: hooks,
		Options:  options,
		Issues:   util.NewSet(),
	}
	return analyzer, nil
}

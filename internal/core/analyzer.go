package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/util"
	"github.com/sirupsen/logrus"
)

// Analyzer git commits analyzer
type Analyzer struct {
	GitHooks   *hook.GitHooks
	Handlers   []Handler
	Info       *hook.Info
	Issues     *util.Set
	Logger     *logrus.Logger
	Options    *config.Options
	Repository *git.Repository
}

// Analyze execute analysis
func (analyzer *Analyzer) Analyze(ctx context.Context, commits []*object.Commit) error {
	if len(analyzer.GitHooks.Hooks) > 0 {
		analyzer.Logger.WithFields(logrus.Fields{
			"correlation_id": util.GetRequestID(ctx),
			"user_id":        util.GetUserID(ctx),
		}).Debugf("%v handlers found and %v hooks found", len(analyzer.Handlers), len(analyzer.GitHooks.Hooks))
		var wg sync.WaitGroup
		for _, commit := range commits {
			wg.Add(1)
			go analyzer.analyze(ctx, &wg, analyzer.GitHooks, commit)
		}
		wg.Wait()
	} else {
		analyzer.Logger.WithFields(logrus.Fields{
			"correlation_id": util.GetRequestID(ctx),
			"user_id":        util.GetUserID(ctx),
		}).Info("There is no hooks in .githooks.yml file")
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
	analyzer.Logger.WithFields(logrus.Fields{
		"correlation_id": util.GetRequestID(ctx),
		"user_id":        util.GetUserID(ctx),
	}).Debugf("Registring handler `%v`", reflect.TypeOf(handler))
	handler.SetLogger(analyzer.Logger)
	if analyzer.Info != nil {
		handler.SetInfo(analyzer.Info)
	}
	handler.SetRepository(analyzer.Repository)
	analyzer.Handlers = append(analyzer.Handlers, handler)
}

// SetHooks set hooks
func (analyzer *Analyzer) SetHooks(hooks *hook.GitHooks) {
	analyzer.GitHooks = hooks
}

// SetInfo set hooks
func (analyzer *Analyzer) SetInfo(info *hook.Info) {
	analyzer.Info = info
}

// SetLogger set logger
func (analyzer *Analyzer) SetLogger(logger *logrus.Logger) {
	analyzer.Logger = logger
}

// SetRepository set repository
func (analyzer *Analyzer) SetRepository(repository *git.Repository) {
	analyzer.Repository = repository
}

func (analyzer *Analyzer) analyze(ctx context.Context, wg *sync.WaitGroup, gitHooks *hook.GitHooks, commit *object.Commit) {
	defer wg.Done()
	issues := make([]issue.Issue, 0)
	for _, hook := range gitHooks.Hooks {
		for _, rule := range hook.Rules {
			for _, handler := range analyzer.Handlers {
				analyzer.Logger.WithFields(logrus.Fields{
					"commit":         commit.Hash.String(),
					"correlation_id": util.GetRequestID(ctx),
					"rule":           rule.Type,
					"user_id":        util.GetUserID(ctx),
				}).Debug("Processing commit rule")
				issuesSlice, _ := handler.Handle(ctx, commit, rule)
				issues = append(issues, issuesSlice...)
			}
		}
	}
	analyzer.Issues.Add(issues)
	statusMessage := util.Colorize(util.Green, config.CheckMark)
	if len(issues) > 0 {
		statusMessage = util.Colorize(util.Red, config.BallotX)
	}
	fmt.Printf("%v · %v · %v\n", commit.Hash.String()[:8], strings.Split(commit.Message, "\n")[0], statusMessage)
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

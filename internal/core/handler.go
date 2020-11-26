package core

import (
	"context"
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/logging"
)

const (
	// HandlerTypeCommits Handler of type 'commits'
	HandlerTypeCommits = "commits"
	// HandlerTypeRefs Handler of type 'refs'
	HandlerTypeRefs = "refs"
)

// Handler hook handler interface
type Handler interface {
	GetRepository() *git.Repository
	GetType() string
	Handle(ctx context.Context, commit *object.Commit, rule *hook.Rule) (issues []issue.Issue, err error)
	SetInfo(info *hook.Info)
	SetLogger(logger logging.Interface)
	SetOptions(options *config.Options)
	SetRepository(repository *git.Repository)
}

// AbstractHandler abstract handler
type AbstractHandler struct {
	Handler
	Info       *hook.Info
	Logger     logging.Interface
	Options    *config.Options
	Repository *git.Repository
}

// GetRepository get git repository
func (handler *AbstractHandler) GetRepository() *git.Repository {
	return handler.Repository
}

// SetInfo set info
func (handler *AbstractHandler) SetInfo(info *hook.Info) {
	handler.Info = info
}

// SetLogger set logger
func (handler *AbstractHandler) SetLogger(logger logging.Interface) {
	handler.Logger = logger
}

// SetOptions set options
func (handler *AbstractHandler) SetOptions(options *config.Options) {
	handler.Options = options
}

// SetRepository set logger
func (handler *AbstractHandler) SetRepository(repository *git.Repository) {
	handler.Repository = repository
}

// CanSkip check if we can skip the rule and condition for given commit
func CanSkip(commit *object.Commit, ruleType hook.HandlerType, conditionType hook.ConditionType) bool {
	var canSkip bool = false
	skipPattern := fmt.Sprintf(`(?i)\[skip[[:space:]]hooks(?:.%s(?:.%s)?)?\]`, ruleType, conditionType)
	if len(regexp.MustCompile(skipPattern).FindStringIndex(commit.Message)) > 0 {
		canSkip = true
	}
	return canSkip
}

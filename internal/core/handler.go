package core

import (
	"context"
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/sirupsen/logrus"
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
	SetLogger(logger *logrus.Logger)
	SetInfo(info *hook.Info)
	SetRepository(repository *git.Repository)
}

// AbstractHandler abstract handler
type AbstractHandler struct {
	Handler
	Info       *hook.Info
	Repository *git.Repository
	Logger     *logrus.Logger
}

// GetRepository get git repository
func (handler *AbstractHandler) GetRepository() *git.Repository {
	return handler.Repository
}

// SetLogger set logger
func (handler *AbstractHandler) SetLogger(logger *logrus.Logger) {
	handler.Logger = logger
}

// SetInfo set logger
func (handler *AbstractHandler) SetInfo(info *hook.Info) {
	handler.Info = info
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

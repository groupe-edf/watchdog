package models

import (
	"github.com/google/uuid"
	"github.com/groupe-edf/watchdog/pkg/query"
)

type Query struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

type Paginator[Model any] struct {
	Items      []Model
	Query      *query.Query
	TotalItems int
}

type Store interface {
	FindCategories(q *query.Query) (Paginator[Category], error)
	GetHealth() error
	GetSettings(q *query.Query) ([]*Setting, error)
	GetWhitelist(q *query.Query) (Paginator[Whitelist], error)
	// Access Tokens
	FindAccessTokens(q *query.Query) (Paginator[AccessToken], error)
	FindAccessToken(name string) (*AccessToken, error)
	SaveAccessToken(token *AccessToken) (*AccessToken, error)
	// Analytics
	Count(container string, q *query.Query) (count int, err error)
	GetAnalytics() ([]AnalyticsData, error)
	GetLeakCountBySeverity() ([]AnalyticsData, error)
	RefreshAnalytics() error
	// Analzes
	DeleteAnalysis(id uuid.UUID) error
	FindAnalyzes(q *query.Query) (Paginator[Analysis], error)
	FindAnalysisByID(id *uuid.UUID) (*Analysis, error)
	SaveAnalysis(analysis *Analysis) (*Analysis, error)
	UpdateAnalysis(analysis *Analysis) (*Analysis, error)
	// Integrations
	AddWebhook(webhook *Webhook) (*Webhook, error)
	DeleteIntegration(id int64) error
	FindIntegrations(q *query.Query) (Paginator[Integration], error)
	FindIntegrationByID(id int64) (*Integration, error)
	FindWebhooks(q *query.Query) (Paginator[Webhook], error)
	SaveIntegration(integration *Integration) (*Integration, error)
	UpdateIntegration(integration *Integration) (*Integration, error)
	// Issues
	FindIssues(q *query.Query) (Paginator[Issue], error)
	SaveIssue(repositoryID *uuid.UUID, analysisID *uuid.UUID, data Issue) error
	// Leaks
	FindLeaks(q *query.Query) (Paginator[Leak], error)
	FindLeakByID(id int64) (*Leak, error)
	SaveLeaks(repositoryID *uuid.UUID, analysisID *uuid.UUID, leaks []Leak) error
	// Policies
	AddPolicyCondition(condition *Condition) (*Condition, error)
	FindPolicies(q *query.Query) (Paginator[Policy], error)
	DeletePolicy(policyID int64) error
	DeletePolicyCondition(policyID int64, conditionID int64) error
	FindPolicyByID(id int64) (*Policy, error)
	SavePolicy(policy *Policy) (*Policy, error)
	TogglePolicy(id int64, enabled bool) error
	UpdatePolicy(policy *Policy) (*Policy, error)
	// Queue
	DeleteJob(job *Job) error
	DoneJob(job *Job)
	Enqueue(job *Job) error
	FindJobs(q *query.Query) (Paginator[*Job], error)
	LockJob(queueName string) (*Job, error)
	SaveJobError(job *Job, message string) error
	// Repositories
	DeleteRepository(id uuid.UUID) error
	FindRepositories(q *query.Query) (Paginator[Repository], error)
	FindRepositoryByID(id *uuid.UUID) (*Repository, error)
	FindRepositoryByURI(uri string) *Repository
	SaveRepository(repository *Repository) (*Repository, error)
	// Rules
	FindRules(q *query.Query) (Paginator[Rule], error)
	SaveRule(rule *Rule) (*Rule, error)
	ToggleRule(ruleID int64, enabled bool) error
	// Users
	DeleteUsers(q *query.Query) error
	FindUserByEmail(email string) (*User, error)
	FindUserById(id *uuid.UUID) (*User, error)
	FindUsers(q *query.Query) (Paginator[User], error)
	SaveUser(user *User) (*User, error)
	SaveOrUpdateUser(user *User) (*User, error)
	UpdatePassword(user *User) (*User, error)
}

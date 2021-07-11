package bolt

import (
	"embed"
	"encoding/json"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/server/query"
	bolt "go.etcd.io/bbolt"
)

const (
	databaseFileName = "watchdog.db"
)

var (
	_ models.Store = &BoltStore{}
	//go:embed migrations/*.json
	migrations embed.FS
	buckets    = []string{
		SettingsBucket,
		UsersBucket,
	}
)

type Migrator struct {
	database *bolt.DB
}

func (migrator *Migrator) Migrate() error {
	entries, err := migrations.ReadDir("migrations")
	if err != nil {
		return err
	}
	for _, migration := range entries {
		file, err := migrations.ReadFile("migrations/" + migration.Name())
		if err != nil {
			return err
		}
		var result map[string]interface{}
		json.Unmarshal([]byte(file), &result)
		for bucketName, data := range result {
			err := migrator.database.Update(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte(bucketName))
				encoded, err := json.Marshal(data)
				if err != nil {
					return err
				}
				return bucket.Put([]byte("SETTINGS"), encoded)
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type BoltStore struct {
	path     string
	database *bolt.DB
}

func (store *BoltStore) Close() error {
	if store.database != nil {
		return store.database.Close()
	}
	return nil
}

func (store *BoltStore) Open() error {
	databasePath := path.Join(store.path, databaseFileName)
	database, err := bolt.Open(databasePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	for _, bucketName := range buckets {
		err := database.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	store.database = database
	return nil
}

func NewBoltStore(options config.BoltOptions) (*BoltStore, error) {
	store := &BoltStore{
		path: options.Path,
	}
	err := store.Open()
	if err != nil {
		return nil, err
	}
	migrator := Migrator{
		database: store.database,
	}
	migrator.Migrate()
	return store, nil
}

func (store *BoltStore) FindCategories(q *query.Query) (models.Paginator[models.Category], error) {
	paginator := models.Paginator[models.Category]{
		Items: make([]models.Category, 0),
		Query: q,
	}
	return paginator, nil
}
func (store *BoltStore) GetHealth() error {
	return nil
}
func (store *BoltStore) GetWhitelist(q *query.Query) (models.Paginator[models.Whitelist], error) {
	paginator := models.Paginator[models.Whitelist]{
		Items: make([]models.Whitelist, 0),
		Query: q,
	}
	return paginator, nil
}

// Analytics
func (store *BoltStore) Count(container string, q *query.Query) (count int, err error) {
	return count, err
}
func (store *BoltStore) GetAnalytics() ([]models.AnalyticsData, error) {
	return nil, nil
}
func (store *BoltStore) GetLeakCountBySeverity() ([]models.AnalyticsData, error) {
	return nil, nil
}
func (store *BoltStore) RefreshAnalytics() error {
	return nil
}

// Analzes
func (store *BoltStore) DeleteAnalysis(id uuid.UUID) error {
	return nil
}
func (store *BoltStore) FindAnalyzes(q *query.Query) (models.Paginator[models.Analysis], error) {
	paginator := models.Paginator[models.Analysis]{
		Items: make([]models.Analysis, 0),
		Query: q,
	}
	return paginator, nil
}
func (store *BoltStore) FindAnalysisByID(id *uuid.UUID) (*models.Analysis, error) {
	return nil, nil
}
func (store *BoltStore) SaveAnalysis(analysis *models.Analysis) (*models.Analysis, error) {
	return nil, nil
}
func (store *BoltStore) UpdateAnalysis(analysis *models.Analysis) (*models.Analysis, error) {
	return nil, nil
}

// Integrations
func (store *BoltStore) FindIntegrations(q *query.Query) (models.Paginator[models.Integration], error) {
	paginator := models.Paginator[models.Integration]{
		Items: make([]models.Integration, 0),
		Query: q,
	}
	return paginator, nil
}
func (store *BoltStore) FindIntegrationByID(id int64) (*models.Integration, error) {
	return nil, nil
}

func (store *BoltStore) SaveIntegration(integration *models.Integration) (*models.Integration, error) {
	return nil, nil
}

func (store *BoltStore) UpdateIntegration(integration *models.Integration) (*models.Integration, error) {
	return nil, nil
}

// Issues
func (store *BoltStore) FindIssues(q *query.Query) (models.Paginator[models.Issue], error) {
	paginator := models.Paginator[models.Issue]{
		Items: make([]models.Issue, 0),
		Query: q,
	}
	return paginator, nil
}
func (store *BoltStore) SaveIssue(repositoryID *uuid.UUID, analysisID *uuid.UUID, data models.Issue) error {
	return nil
}

// Leaks
func (store *BoltStore) FindLeaks(q *query.Query) (models.Paginator[models.Leak], error) {
	paginator := models.Paginator[models.Leak]{
		Items: make([]models.Leak, 0),
		Query: q,
	}
	return paginator, nil
}
func (store *BoltStore) FindLeakByID(id int64) (*models.Leak, error) {
	return nil, nil
}
func (store *BoltStore) SaveLeaks(repositoryID *uuid.UUID, analysisID *uuid.UUID, leaks []models.Leak) error {
	return nil
}

// Queue
func (store *BoltStore) DeleteJob(job *models.Job) error {
	return nil
}
func (store *BoltStore) DoneJob(job *models.Job) {}
func (store *BoltStore) FindJobs(q *query.Query) (models.Paginator[*models.Job], error) {
	paginator := models.Paginator[*models.Job]{
		Items: make([]*models.Job, 0),
		Query: q,
	}
	return paginator, nil
}
func (store *BoltStore) Enqueue(job *models.Job) error {
	return nil
}
func (store *BoltStore) LockJob(queueName string) (*models.Job, error) {
	return nil, nil
}

func (store *BoltStore) SaveJobError(job *models.Job, message string) error {
	return nil
}

// Repositories
func (store *BoltStore) DeleteRepository(id uuid.UUID) error {
	return nil
}
func (store *BoltStore) FindRepositories(q *query.Query) (models.Paginator[models.Repository], error) {
	paginator := models.Paginator[models.Repository]{
		Items: make([]models.Repository, 0),
		Query: q,
	}
	return paginator, nil
}
func (store *BoltStore) FindRepositoryByID(id *uuid.UUID) (*models.Repository, error) {
	return nil, nil
}

func (store *BoltStore) FindRepositoryByURI(uri string) *models.Repository {
	return nil
}
func (store *BoltStore) SaveRepository(repository *models.Repository) (*models.Repository, error) {
	return nil, nil
}

// Users
func (store *BoltStore) FindUserById(id *uuid.UUID) (*models.User, error) {
	return nil, nil
}

func (store *BoltStore) FindUsers(q *query.Query) (models.Paginator[models.User], error) {
	paginator := models.Paginator[models.User]{
		Items: make([]models.User, 0),
		Query: q,
	}
	return paginator, nil
}

func (store *BoltStore) UpdatePassword(user *models.User) (*models.User, error) {
	return nil, nil
}

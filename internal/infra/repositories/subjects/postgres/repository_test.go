package postgres_test

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"

	"github.com/DanilaKorobkov/defi-monitoring/internal/domain"
	pg "github.com/DanilaKorobkov/defi-monitoring/internal/infra/repositories/subjects/postgres"
	"github.com/DanilaKorobkov/defi-monitoring/test/generators"
)

type repositorySuite struct {
	suite.Suite

	db       *sqlx.DB
	tx       *sqlx.Tx
	subjects *pg.SubjectsRepository
}

func TestRepository(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) SetupSuite() {
	projectDir := mustGetProjectDir(getCurrentFile())
	envPath := path.Join(projectDir, "deploy/.env.local")

	err := godotenv.Overload(envPath)
	s.Require().NoError(err)

	config := PostgresConfig{}
	err = env.Parse(&config)
	s.Require().NoError(err)

	db, err := sqlx.Connect("postgres", config.MakeURL())
	s.Require().NoError(err)
	s.db = db
}

func (s *repositorySuite) SetupTest() {
	tx, err := s.db.Beginx()
	s.Require().NoError(err)
	s.tx = tx

	s.subjects = pg.NewSubjectsRepository(tx)
}

func (s *repositorySuite) TearDownTest() {
	err := s.tx.Rollback()
	s.Require().NoError(err)
}

func (s *repositorySuite) TestAdd_AlreadyExists_Override() {
	ctx := context.Background()

	subject := generators.NewSubjectGenerator().Slim().Result()

	err := s.subjects.Add(ctx, subject)
	s.Require().NoError(err)

	subjects, err := s.subjects.GetAll(ctx)
	s.Require().NoError(err)
	s.Require().Equal([]domain.Subject{subject}, subjects)

	overrideSubject := generators.NewSubjectGenerator().Slim().Result()
	overrideSubject.TelegramUserID = subject.TelegramUserID

	err = s.subjects.Add(ctx, overrideSubject)
	s.Require().NoError(err)

	subjects, err = s.subjects.GetAll(ctx)
	s.Require().NoError(err)
	s.Require().Equal([]domain.Subject{overrideSubject}, subjects)
}

func (s *repositorySuite) TestAdd_Success() {
	ctx := context.Background()

	subject := generators.NewSubjectGenerator().Slim().Result()

	err := s.subjects.Add(ctx, subject)
	s.Require().NoError(err)

	subjects, err := s.subjects.GetAll(ctx)
	s.Require().NoError(err)
	s.Require().Equal([]domain.Subject{subject}, subjects)
}

func (s *repositorySuite) TestGetAll_Empty() {
	ctx := context.Background()

	subjects, err := s.subjects.GetAll(ctx)

	s.Require().NoError(err)
	s.Require().Nil(subjects)
}

type PostgresConfig struct {
	PostgresPort     string `env:"POSTGRES_PORT,required"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresDB       string `env:"POSTGRES_DB,required"`
}

func (config PostgresConfig) MakeURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		config.PostgresUser, config.PostgresPassword, config.PostgresPort, config.PostgresDB,
	)
}

func mustGetProjectDir(file string) string {
	for dir := filepath.Dir(file); dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		if filepath.Base(dir) == "internal" {
			return filepath.Dir(dir)
		}
	}

	message := "getProjectDir: cannot find project directory: " + file
	panic(message)
}

func getCurrentFile() string {
	_, filename, _, ok := runtime.Caller(1) //nolint:dogsled // Standard library
	if !ok {
		panic("Could not get current file")
	}
	return filename
}

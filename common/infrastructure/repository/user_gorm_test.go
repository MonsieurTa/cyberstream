package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository *UserGORM
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)
	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	require.NoError(s.T(), err)

	s.DB.Logger.LogMode(logger.Info)

	s.repository = NewUserGORM(s.DB)
}

func (s *Suite) TestCreateUser() {
	model, err := entity.NewUser(entity.CreateUserT{
		FirstName: "ok",
		LastName:  "ok",
		Phone:     "ok",
		Email:     "ok",
		Username:  "ok",
		Password:  "ok",
	})

	require.NoError(s.T(), err)

	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT`)).
		WithArgs(
			model.ID,
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT`)).
		WithArgs(
			model.Credentials.ID,
			model.Credentials.UserID,
			model.Credentials.Username,
			model.Credentials.PasswordHash,
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT`)).
		WithArgs(
			model.PublicInfo.ID,
			model.PublicInfo.UserID,
			model.PublicInfo.FirstName,
			model.PublicInfo.LastName,
			model.PublicInfo.Phone,
			model.PublicInfo.Email,
			model.PublicInfo.PictureURL,
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectCommit()

	_, err = s.repository.Create(model)

	require.NoError(s.T(), err)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

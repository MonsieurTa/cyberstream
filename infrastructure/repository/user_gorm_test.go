package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsieurTa/hypertube/entity"
	"github.com/MonsieurTa/hypertube/infrastructure/model"
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
	model      entity.User
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
	s.model = model.UserModelGORM().(entity.User)
}

func (s *Suite) TestCreateUser() {
	s.model.FillWith(entity.CreateUserT{
		FirstName: "ok",
		LastName:  "ok",
		Phone:     "ok",
		Email:     "ok",
		Username:  "ok",
		Password:  "ok",
	})

	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT`)).
		WithArgs(
			s.model.ID,
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT`)).
		WithArgs(
			s.model.Credentials.ID,
			s.model.Credentials.UserID,
			s.model.Credentials.Username,
			s.model.Credentials.PasswordHash,
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT`)).
		WithArgs(
			s.model.PublicInfo.ID,
			s.model.PublicInfo.UserID,
			s.model.PublicInfo.FirstName,
			s.model.PublicInfo.LastName,
			s.model.PublicInfo.Phone,
			s.model.PublicInfo.Email,
			s.model.PublicInfo.PictureURL,
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectCommit()

	_, err := s.repository.Create(&s.model)

	require.NoError(s.T(), err)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

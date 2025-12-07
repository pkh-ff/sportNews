package repository

import (
	"fmt"
	"os"
	"sportNews/pkg/log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"xorm.io/xorm"
	"xorm.io/xorm/core"
)

func TestMain(m *testing.M) {
	log.InitLogger(false)
	code := m.Run()
	log.CloseLogger()
	os.Exit(code)
}

func TestClose(t *testing.T) {
	eg, mock := newMockEngineGroup(t)

	repo := newMockRepository(eg)

	mock.ExpectClose()

	err := repo.Close()

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)

	repo := newMockRepository(eg)

	mock.ExpectClose().
		WillReturnError(fmt.Errorf("close error"))

	err := repo.Close()

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseWithNilDB(t *testing.T) {
	repo := Repository{
		db:  nil,
		idx: 0,
	}

	err := repo.Close()

	require.NoError(t, err)
}

func newMockEngineGroup(t *testing.T) (*xorm.EngineGroup, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	coreDB := core.FromDB(db)

	engine, err := xorm.NewEngineWithDB("mysql", "sqlmock_db_0", coreDB)
	require.NoError(t, err)

	engineGroup, err := xorm.NewEngineGroup(
		engine,
		[]*xorm.Engine{},
	)
	require.NoError(t, err)

	return engineGroup, mock
}

func newMockRepository(eg *xorm.EngineGroup) *Repository {
	r := New(eg)

	return &r
}

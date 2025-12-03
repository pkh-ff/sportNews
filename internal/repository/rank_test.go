package repository

import (
	"fmt"
	"os"
	"regexp"
	"sportNews/internal/enum"
	"sportNews/internal/model"
	"sportNews/pkg/log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"xorm.io/xorm"
	"xorm.io/xorm/core"
)

const getRankDateSql = "SELECT `data` FROM `sport_rank` WHERE (type = ?) ORDER BY date DESC LIMIT 1"
const insertRankPattern = "INSERT INTO .*sport_rank.*"
const deleteRankPattern = "DELETE FROM .*sport_rank.*"
const getOldestRankSql = "SELECT `id` FROM `sport_rank` WHERE (type = ?) ORDER BY date ASC LIMIT 1"
const checkRankDataExistSql = "SELECT `id`, `type`, `data`, `date`, `create_at` FROM `sport_rank` WHERE (type = ?) AND (date = ?) LIMIT 1"
const getRankDataCountByTypeSql = "SELECT count(*) FROM `sport_rank` WHERE (type = ?)"

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

func TestMain(m *testing.M) {
	log.InitLogger(false)
	code := m.Run()
	log.CloseLogger()
	os.Exit(code)
}

func TestGetRankDate(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	rankType := enum.Test
	expectedData := `{"rank":"some-json-data"}`

	rows := sqlmock.NewRows([]string{"data"}).
		AddRow(expectedData)

	mock.ExpectQuery(regexp.QuoteMeta(getRankDateSql)).
		WithArgs(rankType).
		WillReturnRows(rows)

	got, err := repo.GetRankDate(rankType)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, expectedData, got.Data)
}

func TestGetRankDateWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	rankType := enum.Test

	mock.ExpectQuery(regexp.QuoteMeta(getRankDateSql)).
		WithArgs(rankType).
		WillReturnError(fmt.Errorf("db error"))

	got, err := repo.GetRankDate(rankType)

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Zero(t, got.Data)
}

func TestGetRankDateWithNoData(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	rankType := enum.Test

	rows := sqlmock.NewRows([]string{"data"})

	mock.ExpectQuery(regexp.QuoteMeta(getRankDateSql)).
		WithArgs(rankType).
		WillReturnRows(rows)

	got, err := repo.GetRankDate(rankType)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Zero(t, got.Data)
}

func TestInsertRank(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	rank := model.SportRank{
		Type: enum.Test,
		Data: `{"rank":"some-json-data"}`,
	}

	mock.ExpectExec(insertRankPattern).
		WillReturnResult(sqlmock.NewResult(1, 1)) // lastInsertId, rowsAffected

	err := repo.InsertRank(rank)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertRankWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	rank := model.SportRank{
		Type: enum.Test,
		Data: `{"rank":"some-json-data"}`,
	}

	mock.ExpectExec(insertRankPattern).
		WillReturnError(fmt.Errorf("db error"))

	err := repo.InsertRank(rank)

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOldestRankDataByType(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	rankType := enum.Test
	expectedID := 1

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(expectedID)

	mock.ExpectQuery(regexp.QuoteMeta(getOldestRankSql)).
		WithArgs(rankType).
		WillReturnRows(rows)

	got, err := repo.GetOldestRankDataByType(rankType)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, expectedID, got.Id)
}

func TestGetOldestRankDataByTypeWithNoData(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	rankType := enum.Test

	rows := sqlmock.NewRows([]string{"id"})

	mock.ExpectQuery(regexp.QuoteMeta(getOldestRankSql)).
		WithArgs(rankType).
		WillReturnRows(rows)

	got, err := repo.GetOldestRankDataByType(rankType)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Zero(t, got.Id)
}

func TestGetOldestRankDataByTypeWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	rankType := enum.Test

	mock.ExpectQuery(regexp.QuoteMeta(getOldestRankSql)).
		WithArgs(rankType).
		WillReturnError(fmt.Errorf("db error"))

	got, err := repo.GetOldestRankDataByType(rankType)

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Zero(t, got.Id)
}

func TestGetRankDataCountByType(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	rankType := enum.Test
	var expectedCount int64 = 5

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(expectedCount)

	mock.ExpectQuery(regexp.QuoteMeta(getRankDataCountByTypeSql)).
		WithArgs(rankType).
		WillReturnRows(rows)

	got, err := repo.GetRankDataCountByType(rankType)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, expectedCount, got)
}

func TestGetRankDataCountByTypeWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	rankType := enum.Test

	mock.ExpectQuery(regexp.QuoteMeta(getRankDataCountByTypeSql)).
		WithArgs(rankType).
		WillReturnError(fmt.Errorf("db error"))

	got, err := repo.GetRankDataCountByType(rankType)

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Zero(t, got)
}

func TestCheckRankDataExist(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	date := "2025-01-01"
	rankType := enum.Test

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(int64(1))

	mock.ExpectQuery(regexp.QuoteMeta(checkRankDataExistSql)).
		WithArgs(rankType, date).
		WillReturnRows(rows)

	ok, err := repo.CheckRankDataExist(date, rankType)

	require.NoError(t, err)
	require.True(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckRankDataExistWithNotExists(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	date := "2025-01-01"
	rankType := enum.Test

	rows := sqlmock.NewRows([]string{"id"})

	mock.ExpectQuery(regexp.QuoteMeta(checkRankDataExistSql)).
		WithArgs(rankType, date).
		WillReturnRows(rows)

	ok, err := repo.CheckRankDataExist(date, rankType)

	require.NoError(t, err)
	require.False(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckRankDataExistWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	date := "2025-01-01"
	rankType := enum.Test

	mock.ExpectQuery(regexp.QuoteMeta(checkRankDataExistSql)).
		WithArgs(rankType, date).
		WillReturnError(fmt.Errorf("db error"))

	ok, err := repo.CheckRankDataExist(date, rankType)

	require.Error(t, err)
	require.False(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteRankData(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	rank := model.SportRank{
		Id:   1,
		Type: enum.Test,
		Data: `{"rank":"some-json-data"}`,
	}

	mock.ExpectExec(deleteRankPattern).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ok, err := repo.DeleteRankData(rank)

	require.NoError(t, err)
	require.True(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteRankDataWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	r := New(eg)
	repo := &r

	rank := model.SportRank{
		Id:   1,
		Type: enum.Test,
		Data: `{"rank":"some-json-data"}`,
	}

	mock.ExpectExec(deleteRankPattern).
		WillReturnError(fmt.Errorf("db error"))

	ok, err := repo.DeleteRankData(rank)

	require.Error(t, err)
	require.False(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

package repository

import (
	"fmt"
	"regexp"
	"sportNews/internal/enum"
	"sportNews/internal/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const getRankDateSqlPattern = "SELECT `data` FROM `sport_rank` WHERE (type = ?) ORDER BY date DESC LIMIT 1"
const insertRankPattern = "INSERT INTO .*sport_rank.*"
const deleteRankPattern = "DELETE FROM .*sport_rank.*"
const getOldestRankSqlPattern = "SELECT `id` FROM `sport_rank` WHERE (type = ?) ORDER BY date ASC LIMIT 1"
const checkRankDataExistSqlPattern = "SELECT `id`, `type`, `data`, `date`, `create_at` FROM `sport_rank` WHERE (type = ?) AND (date = ?) LIMIT 1"
const getRankDataCountByTypeSqlPattern = "SELECT count(*) FROM `sport_rank` WHERE (type = ?)"

func TestGetRankDate(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	rankType := enum.Test
	expected := model.SportRank{
		Type: rankType,
		Data: "test",
	}

	rows := sqlmock.NewRows([]string{"data"}).
		AddRow(expected.Data)

	mock.ExpectQuery(regexp.QuoteMeta(getRankDateSqlPattern)).
		WithArgs(rankType).
		WillReturnRows(rows)

	actual, err := repo.GetRankDate(rankType)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, expected.Data, actual.Data)
}

func TestGetRankDateWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	rankType := enum.Test

	mock.ExpectQuery(regexp.QuoteMeta(getRankDateSqlPattern)).
		WithArgs(rankType).
		WillReturnError(fmt.Errorf("db error"))

	actual, err := repo.GetRankDate(rankType)

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Zero(t, actual.Data)
}

func TestGetRankDateWithNoData(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	rankType := enum.Test

	rows := sqlmock.NewRows([]string{"data"})

	mock.ExpectQuery(regexp.QuoteMeta(getRankDateSqlPattern)).
		WithArgs(rankType).
		WillReturnRows(rows)

	actual, err := repo.GetRankDate(rankType)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Zero(t, actual.Data)
}

func TestInsertRank(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	rank := model.SportRank{
		Type: enum.Test,
		Data: "test",
	}

	mock.ExpectExec(insertRankPattern).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.InsertRank(rank)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertRankWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	rank := model.SportRank{
		Type: enum.Test,
		Data: "test",
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

	repo := newMockRepository(eg)

	rankType := enum.Test
	expected := 1

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(expected)

	mock.ExpectQuery(regexp.QuoteMeta(getOldestRankSqlPattern)).
		WithArgs(rankType).
		WillReturnRows(rows)

	actual, err := repo.GetOldestRankDataByType(rankType)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, expected, actual.Id)
}

func TestGetOldestRankDataByTypeWithNoData(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	rankType := enum.Test

	rows := sqlmock.NewRows([]string{"id"})

	mock.ExpectQuery(regexp.QuoteMeta(getOldestRankSqlPattern)).
		WithArgs(rankType).
		WillReturnRows(rows)

	actual, err := repo.GetOldestRankDataByType(rankType)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Zero(t, actual.Id)
}

func TestGetOldestRankDataByTypeWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	rankType := enum.Test

	mock.ExpectQuery(regexp.QuoteMeta(getOldestRankSqlPattern)).
		WithArgs(rankType).
		WillReturnError(fmt.Errorf("db error"))

	actual, err := repo.GetOldestRankDataByType(rankType)

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Zero(t, actual.Id)
}

func TestGetRankDataCountByType(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	rankType := enum.Test
	var expected int64 = 5

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(expected)

	mock.ExpectQuery(regexp.QuoteMeta(getRankDataCountByTypeSqlPattern)).
		WithArgs(rankType).
		WillReturnRows(rows)

	actual, err := repo.GetRankDataCountByType(rankType)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, expected, actual)
}

func TestGetRankDataCountByTypeWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	rankType := enum.Test

	mock.ExpectQuery(regexp.QuoteMeta(getRankDataCountByTypeSqlPattern)).
		WithArgs(rankType).
		WillReturnError(fmt.Errorf("db error"))

	actual, err := repo.GetRankDataCountByType(rankType)

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Zero(t, actual)
}

func TestCheckRankDataExist(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	date := "2025-01-01"
	rankType := enum.Test

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(int64(1))

	mock.ExpectQuery(regexp.QuoteMeta(checkRankDataExistSqlPattern)).
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

	repo := newMockRepository(eg)

	date := "2025-01-01"
	rankType := enum.Test

	rows := sqlmock.NewRows([]string{"id"})

	mock.ExpectQuery(regexp.QuoteMeta(checkRankDataExistSqlPattern)).
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

	repo := newMockRepository(eg)

	date := "2025-01-01"
	rankType := enum.Test

	mock.ExpectQuery(regexp.QuoteMeta(checkRankDataExistSqlPattern)).
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

	repo := newMockRepository(eg)

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

	repo := newMockRepository(eg)

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

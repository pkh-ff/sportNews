package repository

import (
	"fmt"
	"sportNews/internal/enum"
	"sportNews/internal/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	selectNewsPattern = "FROM `news`"
	insertNewsPattern = "INSERT INTO `news`"
	updateNewsPattern = "UPDATE `news`"
)

func TestGetNoCoverNews(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	expected := []model.News{
		{Title: "title1"},
		{Title: "title2"},
	}

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(enum.Enable, "").
		WillReturnRows(newsToRows(expected))

	actual, err := repo.GetNoCoverNews()

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, len(expected), len(actual))
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i].Title, actual[i].Title)
	}
}

func TestGetNoCoverNewsWithDataEmpty(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	rows := sqlmock.NewRows([]string{"data"})

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(enum.Enable, "").
		WillReturnRows(rows)

	actual, err := repo.GetNoCoverNews()

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Empty(t, actual)
}

func TestGetNoCoverNewsWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(enum.Enable, "").
		WillReturnError(fmt.Errorf("db error"))

	actual, err := repo.GetNoCoverNews()

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Zero(t, actual)
}

func TestGetNoCoverCustomNews(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	expected := []model.News{
		{Title: "title1"},
		{Title: "title2"},
	}

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(enum.Enable, "").
		WillReturnRows(newsToRows(expected))

	actual, err := repo.GetNoCoverCustomNews()

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, len(expected), len(actual))
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i].Title, actual[i].Title)
	}
}

func TestGetNoCoverCustomNewsWithEmpty(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	rows := sqlmock.NewRows([]string{"data"})

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(enum.Enable, "").
		WillReturnRows(rows)

	actual, err := repo.GetNoCoverCustomNews()

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Empty(t, actual)
}

func TestGetNoCoverCustomNewsWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(enum.Enable, "").
		WillReturnError(fmt.Errorf("db error"))

	actual, err := repo.GetNoCoverCustomNews()

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Zero(t, actual)
}

func TestGetLastUpdateCoverCustomNews(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	expected := model.News{
		CoverCustom: "CoverCustom",
	}

	rows := sqlmock.NewRows([]string{"cover_custom"}).
		AddRow(expected.CoverCustom)

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(enum.Enable, ""). // status = ?, cover_custom != ?
		WillReturnRows(rows)

	actual, err := repo.GetLastUpdateCoverCustomNews()
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Equal(t, expected, actual)
}

func TestGetLastUpdateCoverCustomNewsWithNoData(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	expected := model.News{}

	rows := sqlmock.NewRows([]string{"cover_custom"})

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(enum.Enable, ""). // status = ?, cover_custom != ?
		WillReturnRows(rows)

	actual, err := repo.GetLastUpdateCoverCustomNews()

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Equal(t, expected, actual)
}

func TestTestGetLastUpdateCoverCustomNewsWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(enum.Enable, "").
		WillReturnError(fmt.Errorf("db error"))

	actual, err := repo.GetLastUpdateCoverCustomNews()

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Zero(t, actual)
}

func TestQueryNewsCount(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	expected := int64(1)
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(expected)

	mock.ExpectQuery(selectNewsPattern).
		WillReturnRows(rows)

	actual, err := repo.QueryNewsCount()

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, expected, actual)
}

func TestQueryNewsCountWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	mock.ExpectQuery(selectNewsPattern).
		WillReturnError(fmt.Errorf("db error"))

	actual, err := repo.QueryNewsCount()

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, int64(0), actual)
}

func TestQueryNewsByPage(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	expected := []model.News{
		{
			Id:          1,
			Title:       "title1",
			Description: "desc1",
			Cover:       "cover1",
			CoverSource: "source1",
			CoverCustom: "custom1",
			PubDate:     time.Now(),
		},
		{
			Id:          2,
			Title:       "title2",
			Description: "desc2",
			Cover:       "cover2",
			CoverSource: "source2",
			CoverCustom: "custom2",
			PubDate:     time.Now(),
		},
	}

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(enum.Enable).
		WillReturnRows(newsPageToRows(expected))

	actual, err := repo.QueryNewsByPage(1, 1)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Equal(t, len(expected), len(actual))
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i].Id, actual[i].Id)
		assert.Equal(t, expected[i].Title, actual[i].Title)
		assert.Equal(t, expected[i].Description, actual[i].Description)
		assert.Equal(t, expected[i].Cover, actual[i].Cover)
		assert.Equal(t, expected[i].CoverSource, actual[i].CoverSource)
		assert.Equal(t, expected[i].CoverCustom, actual[i].CoverCustom)
		assert.True(t, expected[i].PubDate.Equal(actual[i].PubDate))
	}
}

func TestQueryNewsByPageWithDataEmpty(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	rows := sqlmock.NewRows([]string{
		"id",
		"title",
		"description",
		"cover",
		"cover_source",
		"cover_custom",
		"pub_date",
	})

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(enum.Enable).
		WillReturnRows(rows)

	actual, err := repo.QueryNewsByPage(1, 1)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Empty(t, actual)
}

func TestQueryNewsByPageWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(enum.Enable).
		WillReturnError(fmt.Errorf("db error"))

	actual, err := repo.QueryNewsByPage(2, 0)

	require.Error(t, err)
	require.Nil(t, actual)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindNews(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	expected := model.News{
		Title:       "title1",
		Description: "desc1",
		Cover:       "cover1",
		CoverSource: "source1",
		CoverCustom: "custom1",
		Content:     "content1",
		PubDate:     time.Now().Round(0),
	}

	rows := sqlmock.NewRows([]string{
		"title",
		"description",
		"cover",
		"cover_source",
		"cover_custom",
		"content",
		"pub_date",
	}).AddRow(
		expected.Title,
		expected.Description,
		expected.Cover,
		expected.CoverSource,
		expected.CoverCustom,
		expected.Content,
		expected.PubDate)

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(1, enum.Enable).
		WillReturnRows(rows)

	actual, err := repo.FindNews(1)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Equal(t, expected, actual)
}

func TestFindNewsWithDataEmpty(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	rows := sqlmock.NewRows([]string{
		"title",
		"description",
		"cover",
		"cover_source",
		"cover_custom",
		"content",
		"pub_date",
	})

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(1, enum.Enable).
		WillReturnRows(rows)

	actual, err := repo.FindNews(1)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Empty(t, actual)
}

func TestFindNewsWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	mock.ExpectQuery(selectNewsPattern).
		WithArgs(1, enum.Enable).
		WillReturnError(fmt.Errorf("db error"))

	actual, err := repo.FindNews(1)

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Zero(t, actual)
}

func TestInsertNews(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	news := model.News{
		Id:    1,
		Title: "insert title",
	}

	mock.ExpectExec(insertNewsPattern).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.InsertNews(news)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertNewsWithDbError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	news := model.News{
		Id:    1,
		Title: "insert title",
	}

	mock.ExpectExec(insertNewsPattern).
		WillReturnError(fmt.Errorf("db error"))

	err := repo.InsertNews(news)

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateNews(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	news := model.News{
		Id:    1,
		Title: "updated title",
	}

	mock.ExpectExec(updateNewsPattern).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateNews(news)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateNewsWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	news := model.News{
		Id:    1,
		Title: "updated title",
	}

	mock.ExpectExec(updateNewsPattern).
		WillReturnError(fmt.Errorf("db error"))

	err := repo.UpdateNews(news)

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCountByTitle(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	expected := int64(1)
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(expected)

	mock.ExpectQuery(selectNewsPattern).
		WithArgs("", "").
		WillReturnRows(rows)

	actual, err := repo.GetCountByTitle("", "")

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Equal(t, expected, actual)
}

func TestGetCountByTitleWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	mock.ExpectQuery(selectNewsPattern).
		WithArgs("", "").
		WillReturnError(fmt.Errorf("db error"))

	actual, err := repo.GetCountByTitle("", "")

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Equal(t, int64(0), actual)
}

func newsToRows(news []model.News) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"title"})
	for _, n := range news {
		rows.AddRow(n.Title)
	}
	return rows
}

func newsPageToRows(list []model.News) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id",
		"title",
		"description",
		"cover",
		"cover_source",
		"cover_custom",
		"pub_date",
	})
	for _, n := range list {
		rows.AddRow(
			n.Id,
			n.Title,
			n.Description,
			n.Cover,
			n.CoverSource,
			n.CoverCustom,
			n.PubDate,
		)
	}
	return rows
}

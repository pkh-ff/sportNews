package repository

import (
	"fmt"
	"sportNews/internal/enum"
	"sportNews/internal/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const selectVideoPattern = "FROM `video`"

func TestVideoList(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	limit := 2
	expected := make([]model.Video, limit)
	expected = append(expected, model.Video{Title: "title-1", Description: "desc-1", Cover: "cover-1", Link: "link-1"})
	expected = append(expected, model.Video{Title: "title-2", Description: "desc-2", Cover: "cover-2", Link: "link-2"})

	mock.ExpectQuery(selectVideoPattern).
		WithArgs(enum.Enable).
		WillReturnRows(videosToRows(expected))

	actual, err := repo.VideoList(limit)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Equal(t, len(expected), len(actual))

	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i].Title, actual[i].Title)
		assert.Equal(t, expected[i].Description, actual[i].Description)
		assert.Equal(t, expected[i].Cover, actual[i].Cover)
		assert.Equal(t, expected[i].Link, actual[i].Link)
	}
}

func TestVideoListWithNoData(t *testing.T) {
	eq, mock := newMockEngineGroup(t)
	defer eq.Close()

	repo := newMockRepository(eq)

	limit := 0
	expected := make([]model.Video, limit)

	mock.ExpectQuery(selectVideoPattern).
		WithArgs(enum.Enable).
		WillReturnRows(videosToRows(expected))

	actual, err := repo.VideoList(limit)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Equal(t, len(expected), len(actual))
}

func TestVideoListWithDBError(t *testing.T) {
	eg, mock := newMockEngineGroup(t)
	defer eg.Close()

	repo := newMockRepository(eg)

	mock.ExpectQuery(selectVideoPattern).
		WithArgs(enum.Enable).
		WillReturnError(fmt.Errorf("db error"))

	videos, err := repo.VideoList(2)

	require.Error(t, err)
	require.Nil(t, videos)
	require.NoError(t, mock.ExpectationsWereMet())
}

func videosToRows(videos []model.Video) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"title", "description", "cover", "link"})
	for _, v := range videos {
		rows.AddRow(v.Title, v.Description, v.Cover, v.Link)
	}
	return rows
}

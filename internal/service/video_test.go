package service

import (
	"errors"
	"sportNews/internal/model"
	"sportNews/internal/repository/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetVideoList(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	videos := []model.Video{
		{
			Title:       "title1",
			Description: "desc1",
			Cover:       "cover1.jpg",
			Link:        "link1.mp4",
		},
		{
			Title:       "title2",
			Description: "desc2",
			Cover:       "cover2.jpg",
			Link:        "link2.mp4",
		},
	}

	repo := mocks.NewMockVideoRepository(ctrl)
	repo.EXPECT().VideoList(10).Return(videos, nil).Times(1)

	s := &Serv{
		VideoRepo: repo,
	}

	result, err := s.GetVideoList()

	require.NoError(t, err)
	require.Len(t, result, len(videos))
	assert.Equal(t, len(result), len(videos))
	assert.Nil(t, err)
	for i, v := range videos {
		assert.Equal(t, v.Title, result[i].Title)
		assert.Equal(t, v.Description, result[i].Description)
		assert.Contains(t, result[i].Cover, v.Cover)
		assert.Contains(t, result[i].Link, v.Link)
	}
}

func TestGetVideoListWithEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockVideoRepository(ctrl)
	repo.EXPECT().VideoList(10).Return([]model.Video{}, nil).Times(1)

	s := &Serv{
		VideoRepo: repo,
	}

	result, err := s.GetVideoList()

	require.NoError(t, err)
	require.Len(t, result, 0)
	assert.Nil(t, err)
	assert.Empty(t, result)
}

func TestGetVideoListWithDBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockVideoRepository(ctrl)
	repo.EXPECT().VideoList(10).Return([]model.Video{}, errors.New("db error")).Times(1)

	s := &Serv{
		VideoRepo: repo,
	}

	result, err := s.GetVideoList()

	require.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	assert.Nil(t, result)
}

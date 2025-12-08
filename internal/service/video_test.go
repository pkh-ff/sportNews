package service

import (
	"fmt"
	"sportNews/internal/model"
	"sportNews/internal/model/api"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockVideoRepo struct {
	videos []model.Video
	limit  int
	err    error
}

func (m *mockVideoRepo) VideoList(limit int) ([]model.Video, error) {
	m.limit = limit
	return m.videos, m.err
}

func TestGetVideoList(t *testing.T) {
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

	mockRepo := &mockVideoRepo{
		videos: videos,
	}

	s := &Serv{
		VideoRepo: mockRepo,
	}

	result, err := s.GetVideoList()

	require.NoError(t, err)

	resp, ok := result.([]api.VideoResp)
	require.True(t, ok)
	require.Len(t, resp, len(videos))
	assert.Equal(t, 10, mockRepo.limit)
	assert.Nil(t, err)
	for i := range videos {
		assert.Equal(t, videos[i].Title, resp[i].Title)
		assert.Equal(t, videos[i].Description, resp[i].Description)

		assert.Contains(t, resp[i].Cover, videos[i].Cover)
		assert.Contains(t, resp[i].Link, videos[i].Link)
	}
}

func TestGetVideoListWithEmpty(t *testing.T) {
	mockRepo := &mockVideoRepo{
		videos: []model.Video{},
	}

	s := &Serv{
		VideoRepo: mockRepo,
	}

	result, err := s.GetVideoList()

	require.NoError(t, err)

	resp, ok := result.([]api.VideoResp)
	require.True(t, ok)
	require.Len(t, resp, 0)
	assert.Nil(t, err)
	assert.Equal(t, 10, mockRepo.limit)
}

func TestGetVideoListWithDBError(t *testing.T) {
	mockRepo := &mockVideoRepo{
		err: fmt.Errorf("db error"),
	}

	s := &Serv{
		VideoRepo: mockRepo,
	}

	result, err := s.GetVideoList()

	require.Error(t, err)
	assert.Equal(t, 10, mockRepo.limit)
	assert.Nil(t, result)
}

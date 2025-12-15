package service

import (
	"errors"
	"fmt"
	"sportNews/internal/model"
	"sportNews/internal/model/api"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockNewsListRepo struct {
	news        []model.News
	count       int64
	queryErr    error
	countErr    error
	calledLimit int
	calledStart int
	queryCalled bool
	countCalled bool
}

func (m *mockNewsListRepo) QueryNewsByPage(limit, start int) ([]model.News, error) {
	m.queryCalled = true
	m.calledLimit = limit
	m.calledStart = start
	return m.news, m.queryErr
}

func (m *mockNewsListRepo) QueryNewsCount() (int64, error) {
	m.countCalled = true
	return m.count, m.countErr
}

func (m *mockNewsListRepo) FindNews(id int) (model.News, error) {
	panic("")
}

func TestQueryNews(t *testing.T) {
	page := 2
	size := 10
	start := (page - 1) * size

	now := time.Now()

	mockRepo := &mockNewsListRepo{
		news: []model.News{
			{
				Id:          1,
				Title:       "title1",
				Description: "desc1",
				Cover:       "cover1.jpg",
				CoverSource: "source1",
				CoverCustom: "",
				PubDate:     now,
			},
			{
				Id:          2,
				Title:       "title2",
				Description: "desc2",
				Cover:       "cover2.jpg",
				CoverSource: "source2",
				CoverCustom: "cover2.jpg",
				PubDate:     now,
			},
		},
		count: 23,
	}

	s := &Serv{
		NewsRepo: mockRepo,
	}

	resp, err := s.QueryNews(page, size)
	require.NoError(t, err)

	require.True(t, mockRepo.queryCalled)
	require.True(t, mockRepo.countCalled)
	require.Len(t, resp.Records, len(mockRepo.news))
	assert.Equal(t, size, mockRepo.calledLimit)
	assert.Equal(t, start, mockRepo.calledStart)
	assert.Equal(t, mockRepo.count, resp.TotalCount)
	assert.Equal(t, 3, resp.TotalPage)

	for i := range mockRepo.news {
		r := resp.Records[i]
		n := mockRepo.news[i]
		assert.Equal(t, n.Id, r.Id)
		assert.Equal(t, n.Title, r.Title)
		assert.Equal(t, n.Description, r.Description)
		assert.Equal(t, n.CoverSource, r.CoverSource)
		assert.Equal(t, n.Id, r.Id)
		assert.True(t, r.PubDate.Equal(n.PubDate))
		assert.Contains(t, r.Cover, n.Cover)
		assert.Contains(t, r.CoverCustom, n.Cover)
	}
}

func TestQueryNewsWithQueryNewsByPageError(t *testing.T) {
	mockRepo := &mockNewsListRepo{
		queryErr: errors.New("db error"),
	}

	s := &Serv{
		NewsRepo: mockRepo,
	}

	resp, err := s.QueryNews(1, 10)

	require.Error(t, err)
	require.True(t, mockRepo.queryCalled)
	assert.False(t, mockRepo.countCalled)

	require.Len(t, resp.Records, 0)
	assert.Equal(t, int64(0), resp.TotalCount)
	assert.Equal(t, 0, resp.TotalPage)
}

func TestQueryNewsWithQueryNewsCountError(t *testing.T) {
	now := time.Now()

	mockRepo := &mockNewsListRepo{
		news: []model.News{
			{
				Id:          1,
				Title:       "title1",
				Description: "desc1",
				Cover:       "cover1.jpg",
				CoverSource: "source1",
				CoverCustom: "",
				PubDate:     now,
			},
		},
		countErr: errors.New("count error"),
	}

	s := &Serv{
		NewsRepo: mockRepo,
	}

	resp, err := s.QueryNews(1, 10)

	require.Error(t, err)
	require.True(t, mockRepo.queryCalled)
	require.True(t, mockRepo.countCalled)

	require.Len(t, resp.Records, 0)
	assert.Equal(t, int64(0), resp.TotalCount)
	assert.Equal(t, 0, resp.TotalPage)
}

type mockNewsRepo struct {
	news   model.News
	err    error
	called bool
	id     int
}

func (m *mockNewsRepo) QueryNewsByPage(limit, start int) ([]model.News, error) {
	panic("")
}

func (m *mockNewsRepo) QueryNewsCount() (int64, error) {
	panic("")
}

func (m *mockNewsRepo) FindNews(id int) (model.News, error) {
	m.called = true
	m.id = id
	return m.news, m.err
}

func TestFindNewsWithCoverFallbackAndStripURL(t *testing.T) {
	now := time.Now()

	mockRepo := &mockNewsRepo{
		news: model.News{
			Id:          1,
			Title:       "Sample News",
			Description: "Short desc",
			Cover:       "cover1.jpg",
			CoverSource: "source1",
			CoverCustom: "",
			Content:     "This is content with http://example.com and https://foo.bar/path inside.",
			PubDate:     now,
		},
	}

	s := &Serv{
		NewsRepo: mockRepo,
	}

	result, err := s.FindNews(1)
	require.NoError(t, err)
	require.True(t, mockRepo.called)

	assert.Equal(t, 1, mockRepo.id)

	assert.Equal(t, mockRepo.news.Title, result.Title)
	assert.Equal(t, mockRepo.news.Description, result.Description)
	assert.Equal(t, mockRepo.news.CoverSource, result.CoverSource)
	assert.True(t, result.PubDate.Equal(mockRepo.news.PubDate))

	assert.Contains(t, result.Cover, mockRepo.news.Cover)
	assert.Contains(t, result.CoverCustom, mockRepo.news.Cover)

	assert.NotEmpty(t, result.Content)
	assert.NotContains(t, result.Content, "http://")
	assert.NotContains(t, result.Content, "https://")
}

func TestFindNewsWithCustomCover(t *testing.T) {
	now := time.Now()

	mockRepo := &mockNewsRepo{
		news: model.News{
			Id:          2,
			Title:       "Custom Cover News",
			Description: "Desc",
			Cover:       "cover2.jpg",
			CoverSource: "source2",
			CoverCustom: "custom_cover.png",
			Content:     "No URL in this content.",
			PubDate:     now,
		},
	}

	s := &Serv{
		NewsRepo: mockRepo,
	}

	result, err := s.FindNews(2)

	require.NoError(t, err)
	require.True(t, mockRepo.called)

	assert.Equal(t, 2, mockRepo.id)
	assert.Contains(t, result.Cover, mockRepo.news.Cover)
	assert.Contains(t, result.CoverCustom, mockRepo.news.CoverCustom)
	assert.NotContains(t, result.CoverCustom, mockRepo.news.Cover)
}

func TestFindNewsWithRepoError(t *testing.T) {
	mockRepo := &mockNewsRepo{
		err: fmt.Errorf("db error"),
	}

	s := &Serv{
		NewsRepo: mockRepo,
	}

	result, err := s.FindNews(123)

	require.Error(t, err)
	assert.Equal(t, api.NewsDetail{}, result)
	require.True(t, mockRepo.called)
	assert.Equal(t, 123, mockRepo.id)
}

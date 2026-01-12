package service

import (
	"errors"
	"sportNews/internal/model"
	"sportNews/internal/repository/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestQueryNews(t *testing.T) {
	setupAssets(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	news := []model.News{
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
	}

	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().QueryNewsByPage(gomock.Any(), gomock.Any()).
		Return(news, nil).Times(1)
	repo.EXPECT().QueryNewsCount().Return(int64(2), nil).Times(1)

	s := &Serv{
		NewsRepo: repo,
	}

	result, err := s.QueryNews(1, 1)

	require.NoError(t, err)
	require.Len(t, result.Records, len(news))
	require.Nil(t, err)
	require.Equal(t, int(2), result.TotalPage)
	require.Equal(t, int64(2), result.TotalCount)

	for i := range result.Records {
		r := result.Records[i]
		n := news[i]

		require.Equal(t, n.Id, r.Id)
		require.Equal(t, n.Title, r.Title)
		require.Equal(t, n.Description, r.Description)
		require.Equal(t, n.CoverSource, r.CoverSource)
		require.Equal(t, n.Id, r.Id)
		require.True(t, r.PubDate.Equal(n.PubDate))
		require.Contains(t, r.Cover, n.Cover)
		require.Contains(t, r.CoverCustom, n.Cover)
	}
}

func TestQueryNewsWithQueryNewsByPageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().QueryNewsByPage(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("query error")).Times(1)

	s := &Serv{
		NewsRepo: repo,
	}
	result, err := s.QueryNews(1, 1)

	require.Error(t, err)
	require.Equal(t, "query error", err.Error())
	require.Empty(t, result)
}

func TestQueryNewsWithQueryNewsCountError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().QueryNewsByPage(gomock.Any(), gomock.Any()).
		Return([]model.News{}, nil).Times(1)
	repo.EXPECT().QueryNewsCount().
		Return(int64(0), errors.New("query count error")).Times(1)

	s := &Serv{
		NewsRepo: repo,
	}

	result, err := s.QueryNews(1, 1)

	require.Error(t, err)
	require.Equal(t, "query count error", err.Error())
	require.Empty(t, result)
}

func TestFindNewsWithCoverFallbackAndStripURL(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	news := model.News{
		Id:          1,
		Title:       "Sample News",
		Description: "Short desc",
		Cover:       "cover1.jpg",
		CoverSource: "source1",
		CoverCustom: "CoverCustom.jpg",
		Content:     "This is content with http://example.com and https://foo.bar/path inside.",
		PubDate:     now,
	}

	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().FindNews(gomock.Any()).Return(news, nil).Times(1)

	s := &Serv{
		NewsRepo: repo,
	}
	result, err := s.FindNews(1)

	require.NoError(t, err)
	require.Equal(t, news.Title, result.Title)
	require.Equal(t, news.Description, result.Description)
	require.Equal(t, news.CoverSource, result.CoverSource)
	require.True(t, result.PubDate.Equal(news.PubDate))

	require.Equal(t, "https://cdn.test/"+news.Cover, result.Cover)
	require.Equal(t, "https://cdn.test/"+news.CoverCustom, result.CoverCustom)

	require.Equal(t, "This is content with  and  inside.", result.Content)
}

func TestFindNewsWithCustomCover(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()

	news := model.News{
		Id:          1,
		Title:       "Custom Cover News",
		Description: "Desc",
		Cover:       "cover2.jpg",
		CoverSource: "source2",
		CoverCustom: "",
		Content:     "No URL in this content.",
		PubDate:     now,
	}

	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().FindNews(gomock.Any()).Return(news, nil).Times(1)

	s := &Serv{
		NewsRepo: repo,
	}

	result, err := s.FindNews(1)

	require.NoError(t, err)
	require.Equal(t, news.Title, result.Title)
	require.Equal(t, news.Description, result.Description)
	require.Equal(t, "https://cdn.test/"+news.Cover, result.Cover)
	require.Equal(t, news.CoverSource, result.CoverSource)
	require.Equal(t, news.PubDate, result.PubDate)
	require.Equal(t, "https://cdn.test/"+news.Cover, result.CoverCustom)
	require.Equal(t, result.CoverCustom, result.CoverCustom)
	require.Equal(t, result.Content, result.Content)
}

func TestFindNewsWithIdNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().FindNews(gomock.Any()).Return(model.News{Id: 0}, nil).Times(1)

	s := &Serv{
		NewsRepo: repo,
	}
	result, err := s.FindNews(1)
	require.Error(t, err)
	require.Empty(t, result)
}

func TestFindNewsWithRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().FindNews(gomock.Any()).Return(model.News{}, errors.New("find error")).Times(1)

	s := &Serv{
		NewsRepo: repo,
	}
	result, err := s.FindNews(1)

	require.Error(t, err)
	require.Equal(t, "find error", err.Error())
	require.Empty(t, result)
}

package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sportNews/internal/model"
	"sportNews/internal/model/api"
	"sportNews/internal/repository/mocks"
	"sportNews/internal/response"
	"sportNews/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type mockNewsSuccessResp struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data api.NewsPageResp `json:"data"`
}

func TestNewsRouter(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockData := mockNewsResp()

	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().QueryNewsByPage(15, 0).Return(mockData, nil).Times(1)
	repo.EXPECT().QueryNewsCount().Return(int64(31), nil).Times(1)

	serv := &service.Serv{NewsRepo: repo}
	e := newMockTestRouter(t, serv)
	w := mockRequest(t, e, http.MethodGet, "/news")

	var resp mockNewsSuccessResp
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "success", resp.Msg)
	require.Equal(t, int64(31), resp.Data.TotalCount)
	require.Equal(t, 3, resp.Data.TotalPage)
	require.Len(t, resp.Data.Records, 2)
	for i, record := range resp.Data.Records {
		data := mockData[i]
		cover := "https://cdn.test/" + data.Cover

		coverCustom := "https://cdn.test/"
		if record.CoverCustom != "" {
			coverCustom = "https://cdn.test/" + data.CoverCustom
		}

		require.Equal(t, data.Id, record.Id)
		require.Equal(t, data.Title, record.Title)
		require.Equal(t, data.Description, record.Description)
		require.Equal(t, data.CoverSource, record.CoverSource)
		require.Equal(t, cover, record.Cover)
		require.Equal(t, coverCustom, record.CoverCustom)
	}
}

func TestNewsRouterWithQueryNewsByPageError(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().QueryNewsByPage(15, 0).Return(nil, errors.New("error")).Times(1)

	serv := &service.Serv{NewsRepo: repo}
	e := newMockTestRouter(t, serv)
	w := mockRequest(t, e, http.MethodGet, "/news")

	var resp response.HttpResp
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Equal(t, response.NOT_FOUND, resp.Code)
	require.Equal(t, http.StatusText(http.StatusNotFound), resp.Msg)
	require.Nil(t, resp.Data)
}

func TestNewsRouterWithQueryNewsCountError(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().QueryNewsByPage(15, 0).Return([]model.News{}, nil).Times(1)
	repo.EXPECT().QueryNewsCount().Return(int64(31), errors.New("")).Times(1)

	serv := &service.Serv{NewsRepo: repo}
	e := newMockTestRouter(t, serv)
	w := mockRequest(t, e, http.MethodGet, "/news")

	var resp response.HttpResp
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Equal(t, response.NOT_FOUND, resp.Code)
	require.Equal(t, http.StatusText(http.StatusNotFound), resp.Msg)
	require.Nil(t, resp.Data)
}

func mockNewsResp() []model.News {
	return []model.News{
		{
			Id:          1,
			Title:       "t1",
			Description: "d1",
			Cover:       "img/a.jpg",
			CoverSource: "src1",
			CoverCustom: "img/custom/a.jpg",
			PubDate:     time.Now()},
		{
			Id:          2,
			Title:       "t2",
			Description: "d2",
			Cover:       "img/b.jpg",
			CoverSource: "src2",
			CoverCustom: "custom/c2.jpg",
			PubDate:     time.Now(),
		},
	}
}

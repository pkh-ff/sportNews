package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sportNews/internal/model"
	"sportNews/internal/model/api"
	"sportNews/internal/repository/mocks"
	"sportNews/internal/response"
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

	data := mockNewsResp()
	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().QueryNewsByPage(15, 0).Return(data, nil).Times(1)
	repo.EXPECT().QueryNewsCount().Return(int64(31), nil).Times(1)

	e := newMockTestRouter(t, repo)
	w := mockRequest(t, e, http.MethodGet, "/news")

	var resp mockNewsSuccessResp
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "success", resp.Msg)
	require.Equal(t, int64(31), resp.Data.TotalCount)
	require.Equal(t, 3, resp.Data.TotalPage)
	require.Len(t, resp.Data.Records, 2)
	for i, record := range resp.Data.Records {
		data := data[i]

		require.Equal(t, data.Id, record.Id)
		require.Equal(t, data.Title, record.Title)
		require.Equal(t, data.Description, record.Description)
		require.Equal(t, data.CoverSource, record.CoverSource)
		require.Equal(t, getMockAssetsPath(t, data.Cover), record.Cover)
		require.Equal(t, getMockAssetsPath(t, data.CoverCustom), record.CoverCustom)
	}
}

func TestNewsRouterWithQueryNewsByPageError(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().QueryNewsByPage(15, 0).Return(nil, errors.New("error")).Times(1)

	e := newMockTestRouter(t, repo)
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

	e := newMockTestRouter(t, repo)
	w := mockRequest(t, e, http.MethodGet, "/news")

	var resp response.HttpResp
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Equal(t, response.NOT_FOUND, resp.Code)
	require.Equal(t, http.StatusText(http.StatusNotFound), resp.Msg)
	require.Nil(t, resp.Data)
}

type mockNewsDetailSuccessResp struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data api.NewsDetail `json:"data"`
}

func TestNewsByIdRouter(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := mockNewsDetailResp()
	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().FindNews(gomock.Any()).Return(data, nil).Times(1)

	e := newMockTestRouter(t, repo)
	w := mockRequest(t, e, http.MethodGet, "/news/1")

	var resp mockNewsDetailSuccessResp
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "success", resp.Msg)
	require.Equal(t, data.Title, resp.Data.Title)
	require.Equal(t, data.Description, resp.Data.Description)
	require.Equal(t, getMockAssetsPath(t, data.Cover), resp.Data.Cover)
	require.Equal(t, data.CoverSource, resp.Data.CoverSource)
	require.Equal(t, getMockAssetsPath(t, data.CoverCustom), resp.Data.CoverCustom)
}

func TestNewsByIdRouterWithData(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().FindNews(gomock.Any()).Return(model.News{}, nil).Times(1)

	e := newMockTestRouter(t, repo)
	w := mockRequest(t, e, http.MethodGet, "/news/1")

	var resp response.HttpResp
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	actualData, _ := json.Marshal(resp.Data)
	expectedData, _ := json.Marshal(api.NewsDetail{})

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, response.SUCCESS, resp.Code)
	require.Equal(t, "success", resp.Msg)
	require.JSONEq(t, string(expectedData), string(actualData))
}

func TestNewsByIdRouterWithIdIllegal(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockNewsRepository(ctrl)
	e := newMockTestRouter(t, repo)
	w := mockRequest(t, e, http.MethodGet, "/news/aaa")

	var resp response.HttpResp
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, "Missing required parameter error or parameter setting error", resp.Msg)
	require.Equal(t, response.PARAMETER_ERROR, resp.Code)
	require.Nil(t, resp.Data)
}

func TestNewsByIdRouterWithIdIsZero(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockNewsRepository(ctrl)
	e := newMockTestRouter(t, repo)
	w := mockRequest(t, e, http.MethodGet, "/news/0")

	var resp response.HttpResp
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, "Missing required parameter error or parameter setting error", resp.Msg)
	require.Equal(t, response.PARAMETER_ERROR, resp.Code)
	require.Nil(t, resp.Data)
}

func TestNewsByIdRouterWithFindNewsError(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockNewsRepository(ctrl)
	repo.EXPECT().FindNews(gomock.Any()).Return(model.News{}, errors.New("error")).Times(1)

	e := newMockTestRouter(t, repo)
	w := mockRequest(t, e, http.MethodGet, "/news/1")

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
			PubDate:     time.Now(),
		},
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

func mockNewsDetailResp() model.News {
	return model.News{
		Id:          1,
		Title:       "t1",
		Description: "d1",
		Cover:       "img/a.jpg",
		CoverSource: "src1",
		CoverCustom: "img/custom/a.jpg",
		PubDate:     time.Now(),
	}
}

func getMockAssetsPath(t *testing.T, path string) string {
	fullPath := "https://cdn.test/"
	if path != "" {
		fullPath = "https://cdn.test/" + path
	}

	return fullPath
}

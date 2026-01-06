package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sportNews/internal/enum"
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

type mockVideoSuccessResp struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data []api.VideoResp `json:"data"`
}

func TestVideoRouter(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := mockVideoListData(t)
	repo := mocks.NewMockVideoRepository(ctrl)
	repo.EXPECT().VideoList(gomock.Any()).Return(data, nil).Times(1)

	serv := &service.Serv{VideoRepo: repo}
	e := newMockTestRouter(t, serv)
	w := mockRequest(t, e, http.MethodGet, "/video")

	var resp mockVideoSuccessResp
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "success", resp.Msg)
	require.Equal(t, response.SUCCESS, resp.Code)
	require.Equal(t, len(data), len(resp.Data))
	for i, record := range resp.Data {
		d := data[i]

		require.Equal(t, d.Title, record.Title)
		require.Equal(t, d.Description, record.Description)
		require.Equal(t, getMockFullPath(t, d.Cover), record.Cover)
		require.Equal(t, getMockFullPath(t, d.Link), record.Link)
	}
}

func TestVideoRouterWithServiceError(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockVideoRepository(ctrl)
	repo.EXPECT().VideoList(gomock.Any()).Return(nil, errors.New("db error")).Times(1)

	serv := &service.Serv{VideoRepo: repo}
	e := newMockTestRouter(t, serv)
	w := mockRequest(t, e, http.MethodGet, "/video")

	var resp response.HttpResp
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Equal(t, response.NOT_FOUND, resp.Code)
	require.Equal(t, http.StatusText(http.StatusNotFound), resp.Msg)
}

func mockVideoListData(t *testing.T) []model.Video {
	t.Helper()

	now := time.Now()
	return []model.Video{
		{
			Id:          1,
			Title:       "title_1",
			Description: "description_1",
			Cover:       "cover_1",
			Link:        "link_1",
			Status:      enum.Enable,
			CreateAt:    now,
			UpdateAt:    now,
		},
		{
			Id:          2,
			Title:       "title_2",
			Description: "description_2",
			Cover:       "cover_2",
			Link:        "link_2",
			Status:      enum.Disable,
			CreateAt:    now,
			UpdateAt:    now,
		},
	}
}

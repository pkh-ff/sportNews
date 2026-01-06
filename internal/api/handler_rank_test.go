package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sportNews/internal/enum"
	"sportNews/internal/model"
	"sportNews/internal/repository/mocks"
	"sportNews/internal/response"
	"sportNews/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type mockRankSuccessResp struct {
	Code int                `json:"code"`
	Msg  string             `json:"msg"`
	Data []model.RankDetail `json:"data"`
}

func TestRankRouter(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := mockRankData(t)
	repo := mocks.NewMockRankRepository(ctrl)
	repo.EXPECT().GetRankDate(gomock.Any()).Return(data, nil).Times(1)

	serv := &service.Serv{RankRepo: repo}
	e := newMockTestRouter(t, serv)
	w := mockRequest(t, e, http.MethodGet, "/rank/test")

	var resp mockRankSuccessResp
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	var result []model.RankDetail
	_ = json.Unmarshal([]byte(data.Data), &result)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "success", resp.Msg)
	require.Equal(t, response.SUCCESS, resp.Code)

	for i, record := range resp.Data {
		r := result[i]
		require.Equal(t, r.Team, record.Team)
		require.Equal(t, r.Matches, record.Matches)
		require.Equal(t, r.Position, record.Position)
		require.Equal(t, r.Points, record.Points)
		require.Equal(t, getMockFullPath(t, r.Icon), record.Icon)
	}
}

func TestRankRouterWithInvalidType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRankRepository(ctrl)

	serv := &service.Serv{RankRepo: repo}
	e := newMockTestRouter(t, serv)
	w := mockRequest(t, e, http.MethodGet, "/rank/invalidType")

	var resp response.HttpResp
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, "Missing required parameter error or parameter setting error", resp.Msg)
	require.Equal(t, response.PARAMETER_ERROR, resp.Code)
	require.Nil(t, resp.Data)
}

func TestRankRouterWithServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRankRepository(ctrl)
	repo.EXPECT().GetRankDate(gomock.Any()).Return(model.SportRank{}, errors.New("db error")).Times(1)

	serv := &service.Serv{RankRepo: repo}
	e := newMockTestRouter(t, serv)
	w := mockRequest(t, e, http.MethodGet, "/rank/test")

	var resp response.HttpResp
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Equal(t, response.NOT_FOUND, resp.Code)
	require.Equal(t, http.StatusText(http.StatusNotFound), resp.Msg)
}

func mockRankData(t *testing.T) model.SportRank {
	t.Helper()

	dataList := []model.RankDetail{
		{
			Team:     "team_1",
			Position: 1,
			Matches:  1,
			Points:   1,
			Rating:   1,
			Icon:     "icon_1",
		},
		{
			Team:     "team_2",
			Position: 2,
			Matches:  2,
			Points:   2,
			Rating:   2,
			Icon:     "icon_2",
		},
	}

	jsonStr, _ := json.Marshal(dataList)

	return model.SportRank{
		Id:   1,
		Type: enum.Test,
		Data: string(jsonStr),
		Date: time.Now(),
	}
}

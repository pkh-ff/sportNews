package service

import (
	"encoding/json"
	"errors"
	"sportNews/internal/enum"
	"sportNews/internal/model"
	"sportNews/internal/repository/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetRankData(t *testing.T) {
	setupAssets(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	details := []model.RankDetail{
		{
			Team: "New Zealand",
			Icon: "/n8/temp/New_Zealand.png",
		},
		{
			Team: "England",
			Icon: "/n8/temp/England.png",
		},
		{
			Team: "India",
			Icon: "/icons/india.png",
		},
	}

	detailsJSON, _ := json.Marshal(details)
	mockRankData := model.SportRank{
		Type: enum.Test,
		Data: string(detailsJSON),
		Date: time.Now(),
	}

	repo := mocks.NewMockRankRepository(ctrl)
	repo.EXPECT().GetRankDate(gomock.Any()).
		Return(mockRankData, nil).Times(1)

	s := &Serv{
		RankRepo: repo,
	}
	result, err := s.GetRankData(enum.Test)

	require.NoError(t, err)
	require.Len(t, result, len(details))
	for i, r := range result {
		icon := "https://cdn.test" + details[i].Icon
		assert.Equal(t, details[i].Team, r.Team)
		assert.Contains(t, icon, r.Icon)
	}
}

func TestGetRankDataWithEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	detailsJSON, _ := json.Marshal([]model.RankDetail{})

	mockRankData := model.SportRank{
		Type: enum.Test,
		Data: string(detailsJSON),
		Date: time.Now(),
	}

	repo := mocks.NewMockRankRepository(ctrl)
	repo.EXPECT().GetRankDate(gomock.Any()).
		Return(mockRankData, nil).Times(1)

	s := &Serv{
		RankRepo: repo,
	}
	result, err := s.GetRankData(enum.Test)

	require.NoError(t, err)
	require.Empty(t, result)
}

func TestGetRankDataWithDBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRankRepository(ctrl)
	repo.EXPECT().GetRankDate(gomock.Any()).
		Return(model.SportRank{}, errors.New("db error")).Times(1)

	s := &Serv{
		RankRepo: repo,
	}
	result, err := s.GetRankData(enum.Test)

	require.Error(t, err)
	require.Equal(t, "db error", err.Error())
	require.Empty(t, result)
}

func TestTestGetRankDataWithInvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRankData := model.SportRank{
		Type: enum.Test,
		Data: "invalid JSON data",
		Date: time.Now(),
	}

	repo := mocks.NewMockRankRepository(ctrl)
	repo.EXPECT().GetRankDate(gomock.Any()).
		Return(mockRankData, nil).Times(1)

	s := &Serv{
		RankRepo: repo,
	}

	result, err := s.GetRankData(enum.Test)

	require.Error(t, err)
	require.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	require.Empty(t, result)
}

package service

import (
	"encoding/json"
	"fmt"
	"sportNews/internal/enum"
	"sportNews/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRankRepo struct {
	rank     model.SportRank
	err      error
	called   bool
	rankType enum.RankType
}

func (m *mockRankRepo) GetRankDate(t enum.RankType) (model.SportRank, error) {
	m.called = true
	m.rankType = t
	return m.rank, m.err
}

func TestGetRankData(t *testing.T) {
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

	jsonBytes, err := json.Marshal(details)

	mockRepo := &mockRankRepo{
		rank: model.SportRank{
			Data: string(jsonBytes),
		},
	}

	s := &Serv{
		RankRepo: mockRepo,
	}
	result, err := s.GetRankData(enum.Test)

	resp, ok := result.([]model.RankDetail)

	require.NoError(t, err)
	require.True(t, mockRepo.called)
	assert.Equal(t, enum.Test, mockRepo.rankType)

	require.True(t, ok)
	require.Len(t, resp, len(details))
	for i := range resp {
		assert.Equal(t, details[i].Team, resp[i].Team)
		assert.Contains(t, details[i].Icon, resp[i].Icon)
	}
}

func TestGetRankDataWithEmpty(t *testing.T) {
	jsonBytes, err := json.Marshal([]model.RankDetail{})

	mockRepo := &mockRankRepo{
		rank: model.SportRank{
			Data: string(jsonBytes),
		},
	}

	s := &Serv{
		RankRepo: mockRepo,
	}
	result, err := s.GetRankData(enum.Test)

	resp, ok := result.([]model.RankDetail)

	require.NoError(t, err)
	require.True(t, mockRepo.called)
	assert.Equal(t, enum.Test, mockRepo.rankType)

	require.True(t, ok)
	require.Len(t, resp, 0)
}

func TestGetRankDataWithDBError(t *testing.T) {
	mockRepo := &mockRankRepo{
		err: fmt.Errorf("db error"),
	}

	s := &Serv{
		RankRepo: mockRepo,
	}
	result, err := s.GetRankData(enum.Test)

	require.Error(t, err)
	assert.Nil(t, result)
	require.True(t, mockRepo.called)
	assert.Equal(t, enum.Test, mockRepo.rankType)
}

func TestTestGetRankDataWithInvalidJSON(t *testing.T) {
	mockRepo := &mockRankRepo{
		rank: model.SportRank{
			Data: "data",
		},
	}

	s := &Serv{
		RankRepo: mockRepo,
	}

	result, err := s.GetRankData(enum.Test)

	require.Error(t, err)
	assert.Nil(t, result)
}

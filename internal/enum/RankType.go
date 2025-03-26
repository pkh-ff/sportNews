package enum

type RankType string

const (
	Test RankType = "test"
	Odi  RankType = "odi"
	T20  RankType = "t20"
)

func IsRankTypeExist(rank RankType) bool {
	list := []RankType{Test, Odi, T20}

	for _, r := range list {
		if r == rank {
			return true
		}
	}
	return false
}

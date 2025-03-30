package enum

type RankType string

const (
	Test RankType = "test"
	Odi  RankType = "odi"
	T20  RankType = "t20"
)

func IsRankTypeExist(rank RankType) bool {
	for _, r := range RankTypeList() {
		if r == rank {
			return true
		}
	}
	return false
}

func RankTypeList() []RankType {
	return []RankType{Test, Odi, T20}
}

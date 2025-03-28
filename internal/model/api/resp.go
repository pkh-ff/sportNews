package api

type VideoResp struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Cover       string `json:"cover"`
}

package api

type FeedbackReq struct {
	Name    string `json:"name" binding:"required,min=2,max=50"`     // name 最小 2 字元，最大 50 字元
	Email   string `json:"email" binding:"required,email"`           // email 格式驗證
	Subject string `json:"subject" binding:"required,min=3,max=100"` // subject 最小 3 字元，最大 100 字元
	Message string `json:"message" binding:"required,min=3,max=500"` // message 最小 3 字元，最大 500 字元
}

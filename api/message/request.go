package message

type GetRequest struct {
	Topic string `form:"topic" binding:"required"`
	Start int    `form:"start" `
	Limit int    `form:"limit" binding:"required"`
}

package session

type InfoReq struct {
	ClientID string `uri:"client_id" binding:"required"`
}

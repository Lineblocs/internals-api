package model

type Fax struct {
	UserId      int    `json:"user_id"`
	WorkspaceId int    `json:"workspace_id"`
	CallId      int    `json:"call_id"`
	Uri         string `json:"uri"`
	APIId       string `json:"api_id"`
}

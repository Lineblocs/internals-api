package model

type Conference struct {
	Name        string `json:"name"`
	WorkspaceId int    `json:"workspace_id"`
	APIId       string `json:"api_id"`
}

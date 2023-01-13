package model

type Log struct {
	UserId      int     `json:"user_id"`
	WorkspaceId int     `json:"workspace_id"`
	Title       string  `json:"title"`
	Report      string  `json:"report"`
	FlowId      int     `json:"flow_id"`
	Level       *string `json:"report"`
	From        *string `json:"from"`
	To          *string `json:"to"`
}

type LogRoutine struct {
	UserId      int
	WorkspaceId int
	Title       string
	Report      string
	FlowId      int
	Level       string
	From        string
	To          string
}

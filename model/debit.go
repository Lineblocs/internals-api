package model

type Debit struct {
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	Cents        float64 `json:"cents"`
	UserId       int     `json:"user_id"`
	Source       string  `json:"source"`
	ModuleId     int     `json:"module_id"`
	Balance      int     `json:"balance"`
	Status       string  `json:"status"`
	Seconds      float64 `json:"seconds"`
	PlanSnapshot string  `json:"plan_snapshot"`

	//extra request field
	WorkspaceId int    `json:"workspace_id"`
	Number      string `json:"number"`
	Type        string `json:"type"`
}

type DebitAPIParams struct {
	Length          int     `json:"length"`
	RecordingLength float64 `json:"recording_length"`
}
type DebitAPI struct {
	UserId      int            `json:"user_id"`
	WorkspaceId int            `json:"workspace_id"`
	Type        string         `json:"type"`
	Source      string         `json:"source"`
	Params      DebitAPIParams `json:"params"`
}

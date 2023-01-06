package model

type Call struct {
	From         string `json:"from"`
	To           string `json:"to"`
	Status       string `json:"status"`
	Direction    string `json:"direction"`
	Duration     string `json:"duration"`
	UserId       int    `json:"user_id"`
	WorkspaceId  int    `json:"workspace_id"`
	APIId        string `json:"api_id"`
	SourceIp     string `json:"source_ip"`
	ChannelId    string `json:"channel_id"`
	SIPCallId    string `json:"call_id"`
	StartedAt    string `json:"started_at"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	PlanSnapshot string `json:"plan_snapshot"`
}

type CallUpdate struct {
	CallId   int    `json:"call_id"`
	Status   string `json:"status"`
	SourceIp string `json:"source_ip"`
}

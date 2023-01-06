package model

type User struct {
	Id        int
	Username  string
	FirstName string
	LastName  string
	Email     string
}

type PSTNInfo struct {
	IPAddr string `json:"ip_addr"`
	DID    string `json:"did"`
}

type RoutableProvider struct {
	Rate       float64 `json:"rate"`
	DialPrefix string  `json:"dial_prefix"`
	Provider   int     `json:"provider"`
	IPAddress  int     `json:"ip_address"`
}

type CallerIDInfo struct {
	CallerID string `json:"caller_id"`
}

type ExtensionFlowInfo struct {
	FlowId          int               `json:"flow_id"`
	CallerID        string            `json:"caller_id"`
	WorkspaceId     int               `json:"workspace_id"`
	FlowJSON        string            `json:"flow_json"`
	Username        string            `json:"username"`
	Name            string            `json:"name"`
	WorkspaceName   string            `json:"workspace_name"`
	Plan            string            `json:"plan"`
	CreatorId       int               `json:"creator_id"`
	Id              int               `json:"id"`
	APIToken        string            `json:"api_token"`
	APISecret       string            `json:"api_secret"`
	WorkspaceParams *[]WorkspaceParam `json:"workspace_params"`
	FreeTrialStatus string            `json:"free_trial_status"`
}

type CodeFlowInfo struct {
	WorkspaceId     int    `json:"workspace_id"`
	Code            string `json:"code"`
	FlowJSON        string `json:"flow_json"`
	Name            string `json:"name"`
	WorkspaceName   string `json:"workspace_name"`
	Plan            string `json:"plan"`
	CreatorId       int    `json:"creator_id"`
	Id              int    `json:"id"`
	APIToken        string `json:"api_token"`
	APISecret       string `json:"api_secret"`
	FreeTrialStatus string `json:"free_trial_status"`
	FoundCode       bool   `json:"found_code"`
}

type DidNumberInfo struct {
	DidNumber      string `json:"number"`
	DidApiNumber   string `json:"api_number"`
	DidWorkspaceId string `json:"workspace_id"`
	TrunkId        int    `json:"trunk_id"`
}

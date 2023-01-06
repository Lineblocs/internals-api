package model

type Workspace struct {
	Id                  int    `json:"id"`
	CreatorId           int    `json:"creator_id"`
	Name                string `json:"name"`
	BYOEnabled          bool   `json:"byo_enabled"`
	IPWhitelistDisabled bool   `json:"ip_whitelist_disabled"`
	OutboundMacroId     int    `json:"outbound_macro_id"`
	Region              string `json:"region"`
	Plan                string `json:"plan"`
}

type WorkspaceParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type WorkspaceCreatorFullInfo struct {
	Id              int               `json:"id"`
	Workspace       *Workspace        `json:"workspace"`
	WorkspaceName   string            `json:"workspace_name"`
	WorkspaceDomain string            `json:"workspace_domain"`
	WorkspaceId     int               `json:"workspace_id"`
	WorkspaceParams *[]WorkspaceParam `json:"workspace_params"`
	OutboundMacroId int               `json:"outbound_macro_id"`
}

type WorkspaceDIDInfo struct {
	WorkspaceId         int               `json:"workspace_id"`
	Number              string            `json:"number"`
	FlowId              int               `json:"flow_id"`
	FlowJSON            string            `json:"flow_json"`
	WorkspaceName       string            `json:"workspace_name"`
	Name                string            `json:"name"`
	Plan                string            `json:"plan"`
	BYOEnabled          bool              `json:"byo_enabled"`
	IPWhitelistDisabled bool              `json:"ip_whitelist_disabled"`
	OutboundMacroId     int               `json:"outbound_macro_id"`
	CreatorId           int               `json:"creator_id"`
	APIToken            string            `json:"api_token"`
	APISecret           string            `json:"api_secret"`
	WorkspaceParams     *[]WorkspaceParam `json:"workspace_params"`
}

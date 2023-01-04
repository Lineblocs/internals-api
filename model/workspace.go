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

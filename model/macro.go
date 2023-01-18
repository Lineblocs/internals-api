package model

type MacroFunction struct {
	Title        string `json:"title"`
	Code         string `json:"code"`
	CompiledCode string `json:"compiled_code"`
}

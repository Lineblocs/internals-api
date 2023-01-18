package model

type Recording struct {
	Id                 int       `json:"id"`
	UserId             int       `json:"user_id"`
	CallId             *int      `json:"call_id"`
	Size               int       `json:"size"`
	WorkspaceId        int       `json:"workspace_id"`
	APIId              string    `json:"api_id"`
	Tags               *[]string `json:"tags"`
	Trim               bool      `json:"trim"`
	TranscriptionReady bool      `json:"transcription_ready"`
	TranscriptionText  string    `json:"transcription_text"`
	StorageId          string    `json:"storage_id"`
	StorageServerIp    string    `json:"storage_server_ip"`
}

type RecordingTranscription struct {
	RecordingId int    `json:"recording_id"`
	Ready       bool   `json:"ready"`
	Text        string `json:"text"`
}

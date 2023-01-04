package handler

import (
	"lineblocs.com/api/call"
	"lineblocs.com/api/user"
)

type Handler struct {
	callStore call.Store
	userStore user.Store
}

func NewHandler(cs call.Store, us user.Store) *Handler {
	return &Handler{
		callStore: cs,
		userStore: us,
	}
}

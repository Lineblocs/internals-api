package utils

import (
	guuid "github.com/google/uuid"
)

func CreateAPIID(prefix string) string {
	id := guuid.New()
	return prefix + "-" + id.String()
}

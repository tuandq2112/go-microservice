package util

import "github.com/google/uuid"

func GenerateRandomString(length int) string {
	randomUUID := uuid.New()
	return randomUUID.String()[:length]
}

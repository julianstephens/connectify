package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

type cursorPayload struct {
	CreatedAt time.Time `json:"created_at"`
}

func encodeCursor(p cursorPayload) string {
	b, _ := json.Marshal(p)
	return base64.RawURLEncoding.EncodeToString(b)
}

func decodeCursor(s string) (cursorPayload, error) {
	if s == "" {
		return cursorPayload{}, errors.New("empty cursor")
	}
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return cursorPayload{}, err
	}
	var p cursorPayload
	if err := json.Unmarshal(b, &p); err != nil {
		return cursorPayload{}, err
	}
	return p, nil
}

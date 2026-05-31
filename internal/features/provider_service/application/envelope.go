package application

import (
	"encoding/json"
	"fmt"

	sharederrors "ktravels-publicsite-api/internal/shared/errors"
)

type MessageEnvelope struct {
	MessageIdentifier string `json:"MessageIdentifier"`
	Name              string `json:"Name"`
	Data              string `json:"Data"`
}

func UnmarshalEnvelope(body []byte) (*MessageEnvelope, error) {
	var env MessageEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("%w: failed to parse message envelope: %w", sharederrors.ErrInvalidInput, err)
	}
	if env.Data == "" {
		return nil, fmt.Errorf("%w: empty data field", sharederrors.ErrInvalidInput)
	}
	return &env, nil
}

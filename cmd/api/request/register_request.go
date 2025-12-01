package request

import (
	"encoding/json"
)

type RegisterReq struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r RegisterReq) Marshal() ([]byte, error) {
	marshal, err := json.Marshal(r)

	if err != nil {
		return nil, err
	}

	return marshal, nil
}

func (r *RegisterReq) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}

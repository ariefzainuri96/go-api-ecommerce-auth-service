package response

import (
	"encoding/json"
)

// @Model
type LoginResponse struct {
	BaseResponse
	Data LoginData `json:"data"`
}

func (r LoginResponse) Marshal() ([]byte, error) {
	marshal, err := json.Marshal(r)

	if err != nil {
		return nil, err
	}

	return marshal, nil
}

func (r *LoginResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}

// @Model
type LoginData struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

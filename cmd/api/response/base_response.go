package response

import "encoding/json"

type BaseResponse struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
}

func (r BaseResponse) MarshalBaseResponse() ([]byte, error) {
	marshal, err := json.Marshal(r)

	if err != nil {
		return nil, err
	}

	return marshal, nil
}

func (r *BaseResponse) UnmarshalBaseResponse(data []byte) error {
	return json.Unmarshal(data, &r)
}

package responses

type BaseResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

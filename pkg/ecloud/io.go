package ecloud

type Response struct {
	State        string `json:"state"`
	Body         any    `json:"body"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	RequestId    string `json:"requestId"`
}

func (resp *Response) Success() bool {
	return resp.State == "OK"
}

func (resp *Response) GetBody() map[string]any {
	return resp.Body.(map[string]any)
}

func BuildRequestData() map[string]any {
	return nil
}

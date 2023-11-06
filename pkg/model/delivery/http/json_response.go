package http

const (
	success = "success"
	fail    = "fail"
)

type JsonResponse struct {
	Data   any    `json:"data,omitempty"`
	Error  any    `json:"error,omitempty"`
	Errors any    `json:"errors,omitempty"`
	Status string `json:"status"`
}

func NewJsonResponseOnSuccess(data any) JsonResponse {
	return JsonResponse{
		Data: data,
	}
}

func NewJsonResponseOnFailure(err any) JsonResponse {
	switch errorType := err.(type) {
	case error:
		return JsonResponse{
			Error: errorType,
		}
	default:
		return JsonResponse{
			Errors: errorType,
		}
	}
}

func SetStatus(jsonResponse *JsonResponse) {
	if jsonResponse.Data != nil {
		jsonResponse.Status = success
		return
	}
	jsonResponse.Status = fail
}

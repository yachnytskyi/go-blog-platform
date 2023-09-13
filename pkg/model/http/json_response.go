package http

type JsonResponse struct {
	Data   any    `json:"data,omitempty"`
	Error  any    `json:"error,omitempty"`
	Errors any    `json:"errors,omitempty"`
	Status string `json:"status"`
}

func NewJsonResponse(data any) *JsonResponse {
	return &JsonResponse{
		Data: data,
	}
}

func NewJsonResponseWithError(err any) *JsonResponse {
	switch errorType := err.(type) {
	case error:
		return &JsonResponse{
			Error: errorType,
		}
	default:
		return &JsonResponse{
			Errors: errorType,
		}
	}
}

// func NewJsonResponseWithErrors(errors any) *JsonResponse {
// 	return &JsonResponse{
// 		Errors: errors,
// 	}
// }

func SetStatus(jsonResponse *JsonResponse) *JsonResponse {
	if jsonResponse.Data != nil {
		jsonResponse.Status = "success"
		return jsonResponse
	} else if jsonResponse.Error != nil {
		jsonResponse.Status = "fail"
		return jsonResponse
	} else if jsonResponse.Errors != nil {
		jsonResponse.Status = "fail"
		return jsonResponse
	} else {
		jsonResponse.Status = "fail"
		return jsonResponse
	}
}

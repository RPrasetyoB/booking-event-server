package helper

import (
	"booking-event-server/dto"
)

type ResponseWithData struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ResponseWithoutData struct {
	Code    int     `json:"code"`
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Token   *string `json:"token,omitempty"`
}

func Response(params dto.ResponseParams) any {
	var response any
	success := params.StatusCode >= 200 && params.StatusCode <= 299

	if params.Data != nil {
		response = &ResponseWithData{
			Code:    params.StatusCode,
			Success: success,
			Message: params.Message,
			Data:    params.Data,
		}
	} else {
		response = &ResponseWithoutData{
			Code:    params.StatusCode,
			Success: success,
			Message: params.Message,
			Token:   params.Token,
		}
	}

	return response
}

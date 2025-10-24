package context

import (
	"github.com/labstack/echo/v4"
)

type (
	Response struct {
		Code           int         `json:"code"`
		SuccessMessage string      `json:"success_message"`
		ErrorMessage   string      `json:"error_message"`
		Data           interface{} `json:"data"`
		// MetaData       MetaData    `json:"metadata"`
	}

	MetaData struct {
		Limit       int    `json:"limit"`
		TotalCounts int    `json:"total_counts"`
		TotalPages  int    `json:"total_pages"`
		CurrentPage int    `json:"current_page"`
		NextPage    int    `json:"next_page"`
		Sort        string `json:"sort"`
	}
)

func (g *GlobalContext) CreateSuccessResponse(code int, message string, data echo.Map) error {
	resp := &Response{
		Code:           code,
		SuccessMessage: message,
		Data:           data,
		// MetaData:       metaData,
	}
	return g.JSON(code, resp)
}

func (g *GlobalContext) CreateErrorResponse(code int, err error, message string) error {
	resp := &Response{
		Code:         code,
		ErrorMessage: message,
	}

	if err != nil {
		resp.ErrorMessage = err.Error()
	}

	return g.JSON(code, resp)
}

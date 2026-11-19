// Package response writes every API reply in the same envelope:
// { success, message, data, meta, error }.
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    any        `json:"data"`
	Meta    any        `json:"meta"`
	Error   *ErrorInfo `json:"error"`
}

type ErrorInfo struct {
	Code    string            `json:"code"`
	Details map[string]string `json:"details,omitempty"` // field -> problem, for form errors
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Success: true, Message: "Success", Data: data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Body{Success: true, Message: "Created", Data: data})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Fail(c *gin.Context, status int, code, message string, details map[string]string) {
	c.JSON(status, Body{
		Success: false,
		Message: message,
		Error:   &ErrorInfo{Code: code, Details: details},
	})
}

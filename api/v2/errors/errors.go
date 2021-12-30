package errors

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func errorResponse(code int, status string, err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: code,
		StatusText:     status,
		ErrorText:      err.Error(),
	}
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		LogError(err)

		code := 422
		status := "Unprocessable Entity"

		if strings.Contains(err.Error(), "invalid") ||
			strings.Contains(err.Error(), "unexpected") ||
			strings.Contains(err.Error(), "EOF") ||
			strings.Contains(err.Error(), "json") {
			code = 400
			status = "Bad Request"
		}

		if strings.Contains(err.Error(), "Error") {
			code = 500
			status = "Internal Server Error"
		}

		render.Render(w, r, errorResponse(code, status, err))
	}
}

func LogError(err error) {
	if err != nil {
		// Implement real logging
		log.Println(err)
		//debug.PrintStack()
	}
}
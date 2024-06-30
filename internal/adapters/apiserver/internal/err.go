package internal

import "net/http"

// ErrResponse renderer for HTTP failed response
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

var NotFoundErrResponse = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

func NewInternalServerErrResponse(err error) *ErrResponse {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal server error",
		ErrorText:      err.Error(),
	}
}

func NewBadRequestErrResponse(err error) *ErrResponse {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Bad request",
		ErrorText:      err.Error(),
	}
}

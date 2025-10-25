package dtos

type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"error,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}
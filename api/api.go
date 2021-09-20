package api

type (
	Request struct {
		A int64 `json:"a"`
		B int64 `json:"b"`
	}

	Response struct {
		Answer int64 `json:"answer"`
	}

	ErrResponse struct {
		Err     string `json:"err,omitempty"`
		Message string `json:"message,omitempty"`
	}
	DivResponse struct {
		Answer  int64  `json:"answer,omitempty"`
		Err     string `json:"err,omitempty"`
		Message string `json:"message,omitempty"`
	}
	Service interface {
		Add(a, b int64) (int64, error)
		Divide(a, b int64) (int64, error)
	}
)

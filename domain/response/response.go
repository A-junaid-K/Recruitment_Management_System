package response

type ErrResponse struct {
	StatusCode int
	Response   any
	Error      string
}

type SuccessResnpose struct {
	StatusCode int
	Response   any
}

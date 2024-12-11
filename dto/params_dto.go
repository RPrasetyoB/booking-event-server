package dto

type ResponseParams struct {
	Success    bool
	StatusCode int
	Message    string
	Data       any
	Token      *string
}

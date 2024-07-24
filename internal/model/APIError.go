package model

type APIError struct {
	Status int
	Error  error
}

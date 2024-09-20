package jbase

type JsonCommonError struct {
	Code        string `json:"code,omitempty"`
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
}

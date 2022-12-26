package models

type AccountAuthenticationResponse struct {
	Token   string  `json:"token"`
	Profile Account `json:"profile"`
}

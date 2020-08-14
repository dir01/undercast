package undercast

type downloadRequest struct {
	Source string `json:"source"`
}

type loginRequest struct {
	Password string `json:"password"`
}

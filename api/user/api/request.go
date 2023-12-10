package api

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (lr LoginRequest) Validate() bool {
	return lr.Username != "" && len(lr.Username) > 3 &&
		lr.Password != "" && len(lr.Password) >= 8
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (rr RegisterRequest) Validate() bool {
	return rr.Username != "" && len(rr.Username) > 3 &&
		rr.Password != "" && len(rr.Password) >= 8
}

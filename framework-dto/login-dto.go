package frameworkdto

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BearerToken string

type LoginResponseDTO struct {
	Token BearerToken `json:"token"`
}

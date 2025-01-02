package schemas

// REQUEST BODY SCHEMAS
type RegisterUser struct {
	FirstName      string `json:"first_name" validate:"required,max=50" example:"Donald"`
	LastName       string `json:"last_name" validate:"required,max=50" example:"Trump"`
	Email          string `json:"email" validate:"required,min=5,email" example:"donaldtrump47th@gmail.com"`
	Password       string `json:"password" validate:"required,min=8,max=50" example:"!2x8w6?0gO94_4,v"`
	TermsAgreement bool   `json:"terms_agreement" validate:"eq=true"`
}

type EmailRequestSchema struct {
	Email string `json:"email" validate:"required,min=5,email" example:"donaldtrump47th@gmail.com"`
}

type VerifyEmailRequestSchema struct {
	EmailRequestSchema
	Otp uint32 `json:"otp" validate:"required" example:"123456"`
}

type SetNewPasswordSchema struct {
	VerifyEmailRequestSchema
	Password string `json:"password" validate:"required,min=8,max=50" example:"H14l@6c$9W{ED?18"`
}

type LoginSchema struct {
	Email    string `json:"email" validate:"required,email" example:"donaldtrump47th@gmail.com"`
	Password string `json:"password" validate:"required" example:"password"`
}

type RefreshTokenSchema struct {
	Refresh string `json:"refresh" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"`
}

// RESPONSE BODY SCHEMAS
type RegisterResponseSchema struct {
	ResponseSchema
	Data EmailRequestSchema `json:"data"`
}

type TokensResponseSchema struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type LoginResponseSchema struct {
	ResponseSchema
	Data TokensResponseSchema `json:"data"`
}

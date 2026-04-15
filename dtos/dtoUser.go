package dtos

type UserResponseDto struct {
	Id        int    `json:"Id"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	CI        int    `json:"CI"`
	DateBirth string `json:"date_birth"`
}

type UserCreateDto struct {
	Token     string `json:"token" validate:"required"`
	Name      string `json:"name" validate:"omitempty,max=50,is_name"`
	LastName  string `json:"last_name" validate:"omitempty,max=50,is_name"`
	CI        int    `json:"CI" validate:"required,gt=11111,lte=999999999"`
	DateBirth string `json:"date_birth" validate:"omitempty,datetime=02/01/2006,age_gte_16"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=7,secure_password"`
}

type UserUpdateDto struct {
	Name      string `json:"name" validate:"omitempty,max=50,is_name"`
	LastName  string `json:"last_name" validate:"omitempty,max=50,is_name"`
	CI        int    `json:"CI" validate:"omitempty,gt=11111,lte=999999999"`
	DateBirth string `json:"date_birth" validate:"omitempty,datetime=02/01/2006,age_gte_16"`
}

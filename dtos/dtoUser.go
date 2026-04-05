package dtos

type UserResponseDto struct{
	Id		int		`json:"Id"`
	Name	string	`json:"name"`
	LastName string `json:"last_name"`
	CI		int		`json:"CI"`
	DateBirth string `json:"date_birth"`
	Email	string	`json:"email"`
}

type UserCreateDto struct{
	Name	string	`json:"name" validate:"required,max=50,is_name"`
	LastName string `json:"last_name" validate:"required,max=50,is_name"`
	CI		int		`json:"CI" validate:"required,gt=11111,lte=999999999"`
	Rol		string	`json:"rol" validate:"required,oneof=DEVELOPER USER"`
	DateBirth string `json:"date_birth" validate:"required,datetime=02/01/2006"`
	Email	string	`json:"email" validate:"required,email"`
	Password string	`json:"password" validate:"required,min=7,secure_password"`
}

type UserUpdateDto struct{
	Name	string	`json:"name" validate:"omitempty,max=50,is_name"`
	LastName string `json:"last_name" validate:"omitempty,max=50,is_name"`
	CI		int		`json:"CI" validate:"omitempty,gt=11111,lte=999999999"`
	DateBirth string `json:"date_birth" validate:"omitempty,datetime=02/01/2006"`
}

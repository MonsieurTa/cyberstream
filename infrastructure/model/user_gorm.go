package model

import (
	"github.com/MonsieurTa/hypertube/entity"
)

func UserModelGORM() interface{} {
	params := []Params{
		{"ID", `gorm:"column:id;type:uuid;not null"`},
		{"CreatedAt", `gorm:"column:created_at"`},
	}
	model := entity.User{}
	return GenerateModelGORM(params, model)
}

func CredentialsModelGORM() interface{} {
	params := []Params{
		{"ID", `gorm:"column:id;type:uuid;not null"`},
		{"UserID", `gorm:"column:user_id"`},
		{"Username", `gorm:"column:username;unique"`},
		{"PasswordHash", `gorm:"column:password_hash"`},
		{"UpdatedAt", `gorm:"column:updated_at"`},
	}
	model := entity.Credentials{}
	return GenerateModelGORM(params, model)
}

func PublicInfoModelGORM() interface{} {
	params := []Params{
		{"ID", `gorm:"column:id;type:uuid;not null"`},
		{"UserID", `gorm:"column:user_id"`},
		{"FirstName", `gorm:"column:first_name"`},
		{"LastName", `gorm:"column:last_name"`},
		{"Phone", `gorm:"column:phone"`},
		{"Email", `gorm:"column:email;unique"`},
		{"PictureURL", `gorm:"column:picture_url"`},
		{"UpdatedAt", `gorm:"column:updated_at"`},
	}
	model := entity.PublicInfo{}
	return GenerateModelGORM(params, model)
}

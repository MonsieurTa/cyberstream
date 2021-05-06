package validator

import (
	"errors"

	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/MonsieurTa/hypertube/server-api/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRegistrationValidator struct {
	User struct {
		PublicInfo struct {
			FirstName string `json:"firstname" binding:"required,alpha"`
			LastName  string `json:"lastname" binding:"required,alpha"`
			Phone     string `json:"phone" binding:"required,phone_fr"`
			Email     string `json:"email" binding:"required,email"`
		} `json:"public_info"`
		Credentials struct {
			Username string `json:"username" binding:"required,alphanum,min=4,max=255"`
			Password string `json:"password" binding:"required,alphanum,min=8,max=255"`
		} `json:"credentials"`
	} `json:"user"`

	output entity.User `json:"-"`
}

func NewUserRegistrationValidator() *UserRegistrationValidator {
	return &UserRegistrationValidator{}
}

func (v *UserRegistrationValidator) Validate(c *gin.Context) error {
	err := common.Bind(c, v)
	if err != nil {
		return err
	}
	input := entity.CreateUserT{
		FirstName: v.User.PublicInfo.FirstName,
		LastName:  v.User.PublicInfo.LastName,
		Phone:     v.User.PublicInfo.Phone,
		Email:     v.User.PublicInfo.Email,
		Username:  v.User.Credentials.Username,
		Password:  v.User.Credentials.Password,
	}
	output, err := entity.NewUser(input)
	if err != nil {
		return err
	}
	v.output = *output
	return nil
}

func (v UserRegistrationValidator) Value() entity.User {
	return v.output
}

type UserUpdateValidator struct {
	Update struct {
		Username   string `json:"username,omitempty"`
		Email      string `json:"email,omitempty"`
		Password   string `json:"password,omitempty"`
		PictureURL string `json:"picture_url,omitempty"`
	} `json:"update"`

	output UserUpdateOutput `json:"-"`
}

type UserUpdateOutput struct {
	UserID     uuid.UUID
	Username   string
	Email      string
	Password   string
	PictureURL string
}

func NewUserUpdateValidator() *UserUpdateValidator {
	return &UserUpdateValidator{}
}

func (v *UserUpdateValidator) Validate(c *gin.Context) error {
	err := common.Bind(c, v)
	if err != nil {
		return err
	}
	token, ok := c.Get("token")
	if !ok {
		return errors.New("unauthorized")
	}

	userID, err := uuid.Parse(token.(*jwt.StandardClaims).Audience)
	if err != nil {
		return err
	}

	v.output.UserID = userID
	v.output.Username = v.Update.Username
	v.output.Password = v.Update.Password
	v.output.Email = v.Update.Email
	v.output.PictureURL = v.Update.PictureURL
	return nil
}

func (v *UserUpdateValidator) Value() UserUpdateOutput {
	return v.output
}

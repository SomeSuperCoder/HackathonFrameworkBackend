package validators

import (
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserValidator struct {
	av       *AccessValidator
	userAuth *models.User
	targetID bson.ObjectID
}

func NewUserValidator(userAuth *models.User, targetID bson.ObjectID) *UserValidator {
	uv := &UserValidator{
		userAuth: userAuth,
		av:       NewAccessValidator(userAuth),
		targetID: targetID,
	}

	// Access
	uv.av.validator.RegisterValidation("self", uv.validateSelf)

	return uv
}

func (uv *UserValidator) validateSelf(fl validator.FieldLevel) bool {
	return uv.av.isAdmin() || uv.userAuth.ID == uv.targetID
}

func (uv *UserValidator) ValidateRequest(r any) error {
	return uv.av.validator.Struct(r)
}

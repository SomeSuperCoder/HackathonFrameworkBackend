package validators

import (
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/go-playground/validator/v10"
)

type AccessValidator struct {
	validator *validator.Validate
	userAuth  *models.User
}

func NewAccessValidator(userAuth *models.User) *AccessValidator {
	v := validator.New()
	av := &AccessValidator{
		validator: v,
		userAuth:  userAuth,
	}
	_ = v.RegisterValidation("admin", av.validateAdmin)
	_ = v.RegisterValidation("judge", av.validateJudge)

	return av
}

func (av *AccessValidator) validateAdmin(fl validator.FieldLevel) bool {
	return av.isAdmin()
}

func (av *AccessValidator) validateJudge(fl validator.FieldLevel) bool {
	return av.isAdmin() || av.userAuth.Role == models.Judge
}

func (av *AccessValidator) isAdmin() bool {
	return av.userAuth.Role == models.Admin
}

func (av *AccessValidator) ValidateRequest(r any) error {
	return av.validator.Struct(r)
}

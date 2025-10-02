package validators

import (
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/go-playground/validator/v10"
)

type TeamValidator struct {
	av       *AccessValidator
	userAuth *models.User
	team     *models.Team
}

func NewTeamValidator(userAuth *models.User, team *models.Team) *TeamValidator {
	tv := &TeamValidator{
		userAuth: userAuth,
		av:       NewAccessValidator(userAuth),
		team:     team,
	}

	// Access
	tv.av.validator.RegisterValidation("owner", tv.validateOwner)
	// Other
	tv.av.validator.RegisterValidation("grades", models.ValidateGrades)

	return tv
}

func (tv *TeamValidator) validateOwner(fl validator.FieldLevel) bool {
	return tv.av.isAdmin() || tv.team.Leader == tv.userAuth.ID
}

func (tv *TeamValidator) ValidateRequest(r any) error {
	return tv.av.validator.Struct(r)
}

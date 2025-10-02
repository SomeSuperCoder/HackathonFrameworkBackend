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
	tv.av.validator.RegisterValidation("grades", tv.validateGrades)

	return tv
}

func (tv *TeamValidator) validateOwner(fl validator.FieldLevel) bool {
	return tv.av.isAdmin() || tv.team.Leader == tv.userAuth.ID
}

func (tv *TeamValidator) validateGrades(fl validator.FieldLevel) bool {
	// Get the field value and check if it's the correct type
	if nestedMap, ok := fl.Field().Interface().(models.Grades); ok {
		// Iterate through the outer map
		for k, innerMap := range nestedMap {
			// Check judge access
			if k != tv.userAuth.ID {
				return false
			}

			// Iterate through each inner map
			for _, value := range innerMap {
				// Check if the integer value is within the required range
				if value > 10 {
					return false
				}
			}
		}
		return true // All values are valid
	}
	return false // Type assertion failed
}

func (tv *TeamValidator) ValidateRequest(r any) error {
	return tv.av.validator.Struct(r)
}

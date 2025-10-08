package validators

type Validator interface {
	ValidateRequest(r any) error
}

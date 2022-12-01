package validator

type IValidator interface {
	Validate(httpMethod string) Message
}

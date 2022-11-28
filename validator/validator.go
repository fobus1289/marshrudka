package validator

type IValidator interface {
	Validate() MessageMapResult
}

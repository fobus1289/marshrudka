package validator

type String struct {
	*Array[string]
}

func StringValidator(key string, t string) *String {
	return &String{
		Array: &Array[string]{
			&AnyType[string, int]{
				value: t,
				key:   key,
			},
		},
	}
}

func (str *String) Min(min int, message string) *String {
	str.Array.Min(min, message)
	return str
}

func (str *String) Max(max int, message string) *String {
	str.Array.Max(max, "max")
	return str
}

func (str *String) Required(message string) *String {
	str.Array.Required(message)
	return str
}

func (str *String) Lenght(min, max int, message string) *String {
	str.Array.Lenght(min, max, message)
	return str
}

func (str *String) Only(message string, strs ...string) *String {

	for _, s := range strs {
		if s == str.value {
			return str
		}
	}

	// str.Array.AddMessage(true, "only", 1)
	return str
}

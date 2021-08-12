package stream

type String string
type Strings []string

func (s String) Filter(fn func(elm int32) bool) String {
	var result String
	for _, e := range s {
		if fn(e) {
			result += String(e)
		}
	}
	return result
}

func (s String) Get() string {
	return string(s)
}

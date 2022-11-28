package validator

type IMessage interface {
	Message() map[string]any
	IsValid() bool
	Key() string
}

type MessageMapResult map[string]any

func Build(messages ...IMessage) MessageMapResult {

	var result = MessageMapResult{}

	for _, message := range messages {
		if message.IsValid() || len(message.Message()) == 0 {
			continue
		}
		result[message.Key()] = message.Message()
	}

	return result
}

func (mr MessageMapResult) ErrorCount() int {
	return len(mr)
}

func (mr MessageMapResult) IsValid() bool {
	return mr.ErrorCount() == 0
}

package validator

type IMessage interface {
	ErrorMessage() Message
}

type MessageMapResult map[string]any

func Build(messages ...IMessage) Message {

	var result = Message{}

	for _, message := range messages {

		mgs := message.ErrorMessage()

		if mgs.Len() == 0 {
			continue
		}

		for k, v := range mgs {
			result[k] = v
		}

	}

	return result
}

func (mr MessageMapResult) ErrorCount() int {
	return len(mr)
}

func (mr MessageMapResult) IsValid() bool {
	return mr.ErrorCount() == 0
}

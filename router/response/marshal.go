package response

import (
	"encoding/json"
	"encoding/xml"
)

var _emptyBody = make([]byte, 0, 0)

func jsonSend(data interface{}) []byte {

	if data == nil {
		return _emptyBody
	}

	if sendBody, err := json.Marshal(data); err == nil {
		return sendBody
	}

	return _emptyBody
}

func xmlSend(data interface{}) []byte {

	if data == nil {
		return _emptyBody
	}

	if sendBody, err := xml.Marshal(data); err == nil {
		return sendBody
	}

	return _emptyBody
}

func fileSend(data interface{}) []byte {
	return _emptyBody
}

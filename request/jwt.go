package request

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	invalidSignatureError = errors.New("invalid signature")
	tokenExpiredError     = errors.New("token expired")
	jwtUserError          = errors.New("jwt user can be null")
	jwtHeader             = []byte(`{"alg":"HS256","typ":"JWT"}`)
)

type IJwtUser interface {
	Build(expired, iat time.Duration) IJwtUser
	Expired() time.Duration
	Out(token string) any
}

type Jwt struct {
	Secret  []byte
	Expired time.Duration
}

func (j *Jwt) Decode(token string, user IJwtUser) error {

	if user == nil {
		return jwtUserError
	}

	if token == "" || strings.Count(token, ".") != 2 {
		return invalidSignatureError
	}

	data := strings.Split(token, ".")

	var (
		header    = data[0]
		payload   = data[1]
		signature = data[2]
	)

	if !j.Valid([]byte(header+payload), []byte(signature)) {
		return invalidSignatureError
	}

	bson, err := base64.StdEncoding.DecodeString(payload)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(bson, user); err != nil {
		return err
	}

	return nil
}

func (j *Jwt) DecodeWithExpired(token string, user IJwtUser) error {

	if err := j.Decode(token, user); err != nil {
		return err
	}

	if time.Now().Unix() > int64(user.Expired()) {
		return tokenExpiredError
	}

	return nil
}

func (j *Jwt) Encode(auth IJwtUser) (string, error) {

	exp := time.Duration(time.Now().Add(j.Expired * time.Minute).Unix())
	iat := time.Duration(time.Now().Unix())

	data, err := json.Marshal(auth.Build(exp, iat))

	if err != nil {
		return "", err
	}

	header := base64.StdEncoding.EncodeToString(jwtHeader)

	payload := base64.StdEncoding.EncodeToString(data)

	signature := base64.StdEncoding.EncodeToString(
		j.NewHash([]byte(header + payload)),
	)

	return strings.Join([]string{header, payload, signature}, "."), nil
}

func (j *Jwt) Valid(message, messageMAC []byte) bool {

	h := hmac.New(sha256.New, []byte(j.Secret))

	result, _ := base64.StdEncoding.DecodeString(string(messageMAC))

	return hmac.Equal(h.Sum(message), result)
}

func (j *Jwt) NewHash(signature []byte) []byte {
	mac := hmac.New(sha256.New, j.Secret)
	return mac.Sum(signature)
}

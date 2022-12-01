package validator

import (
	"fmt"
	"regexp"
	"strings"
)

type IStringValidator interface {
	Min(min int, message string) IStringValidator
	Max(max int, message string) IStringValidator
	Lenght(min, max int, message string) IStringValidator
	Only(message string, strs ...string) IStringValidator
	Omit(message string, strs ...string) IStringValidator
	OnlyIgnoreCase(message string, strs ...string) IStringValidator
	OmitIgnoreCase(message string, strs ...string) IStringValidator
	Email(message string) IStringValidator
	Regular(message, pattern string) IStringValidator
	Equals(message, s string) IStringValidator
	EqualsIgnoreCase(message, s string) IStringValidator
	Prefix(message, s string) IStringValidator
	Mask(mask, message string) IStringValidator
	Options(optional bool) IStringValidator
	IMessage
}

type String struct {
	Key      string
	Value    *string
	Optional bool
	Message
}

func StringValidator(key string, value *string) IStringValidator {
	return &String{
		Key:      key,
		Value:    value,
		Optional: false,
		Message:  Message{},
	}
}

func (str *String) Options(optional bool) IStringValidator {
	str.Optional = optional
	return str
}

func (str *String) Min(min int, message string) IStringValidator {

	if str.Optional && str.Value == nil {
		return str
	}

	if len(*str.Value) < min {
		str.Add(str.Key, "min", message, str.Value)
	}
	return str
}

func (str *String) Max(max int, message string) IStringValidator {
	var value string
	{
		if str.Value != nil {
			value = *str.Value
		}
	}

	if str.Optional && str.Value == nil {
		return str
	}

	if len(value) > max {
		str.Add(str.Key, "max", message, str.Value)
	}

	return str
}

func (str *String) Lenght(min, max int, message string) IStringValidator {
	var value string
	{
		if str.Value != nil {
			value = *str.Value
		}
	}

	if str.Optional && str.Value == nil {
		return str
	}

	l := len(value)

	if l < min || l > max {
		str.Add(str.Key, "lenght", message, str.Value)
	}

	return str
}

func (str *String) Only(message string, strs ...string) IStringValidator {
	if len(strs) == 0 {
		return str
	}

	var value string
	{
		if str.Value != nil {
			value = *str.Value
		}
	}

	if str.Optional && str.Value == nil {
		return str
	}

	for _, s := range strs {
		if s == value {
			return str
		}
	}

	str.Add(str.Key, "only", message, str.Value)

	return str
}

func (str *String) Omit(message string, strs ...string) IStringValidator {
	if len(strs) == 0 {
		return str
	}

	var value string
	{
		if str.Value != nil {
			value = *str.Value
		}
	}

	if str.Optional && str.Value == nil {
		return str
	}

	for _, s := range strs {
		if s != value {
			return str
		}
	}

	str.Add(str.Key, "omit", message, str.Value)

	return str
}

func (str *String) OnlyIgnoreCase(message string, strs ...string) IStringValidator {

	if len(strs) == 0 {
		return str
	}

	var value string
	{
		if str.Value != nil {
			value = *str.Value
		}
	}

	if str.Optional && str.Value == nil {
		return str
	}

	for _, s := range strs {
		if strings.EqualFold(value, s) {
			return str
		}
	}

	str.Add(str.Key, "onlyIgnoreCase", message, str.Value)

	return str
}

func (str *String) OmitIgnoreCase(message string, strs ...string) IStringValidator {

	if len(strs) == 0 {
		return str
	}

	var value string
	{
		if str.Value != nil {
			value = *str.Value
		}
	}

	if str.Optional && str.Value == nil {
		return str
	}

	for _, s := range strs {
		if !strings.EqualFold(value, s) {
			return str
		}
	}

	str.Add(str.Key, "omitIgnoreCase", message, str.Value)

	return str
}

func (str *String) Email(message string) IStringValidator {

	var value string
	{
		if str.Value != nil {
			value = *str.Value
		}
	}

	if str.Optional && str.Value == nil {
		return str
	}

	if value == "" || !isEmailValid(value) {
		str.Add(str.Key, "email", message, str.Value)
	}

	return str
}

func (str *String) Regular(message, pattern string) IStringValidator {

	var value string
	{
		if str.Value != nil {
			value = *str.Value
		}
	}

	if str.Optional && str.Value == nil {
		return str
	}

	reg, err := regexp.Compile(pattern)

	if err != nil {
		str.Add(str.Key, "regular", err.Error(), pattern)
		return str
	}

	if !reg.MatchString(value) {
		str.Add(str.Key, "regular", message, str.Value)
	}

	return str
}

func (str *String) Equals(message, s string) IStringValidator {

	var value string
	{
		if str.Value != nil {
			value = *str.Value
		}
	}

	if str.Optional && str.Value == nil {
		return str
	}

	if value != s {
		str.Add(str.Key, "equals", message, str.Value)
	}

	return str
}

func (str *String) EqualsIgnoreCase(message, s string) IStringValidator {

	var value string
	{
		if str.Value != nil {
			value = *str.Value
		}
	}

	if str.Optional && str.Value == nil {
		return str
	}

	if !strings.EqualFold(value, s) {
		str.Add(str.Key, "equalsIgnoreCase", message, str.Value)
	}

	return str
}

func (str *String) Prefix(message, s string) IStringValidator {

	var value string
	{
		if str.Value != nil {
			value = *str.Value
		}
	}

	if str.Optional && str.Value == nil {
		return str
	}

	if !strings.HasPrefix(value, s) {
		str.Add(str.Key, "prefix", message, str.Value)
	}

	return str
}

func (str *String) Mask(mask, message string) IStringValidator {

	mask = strings.ReplaceAll(mask, "#", `\w`)

	reg, err := regexp.Compile(fmt.Sprintf("^(%s)$", mask))

	if err != nil {
		str.Add(str.Key, "mask", err.Error(), mask)
		return str
	}

	var value string
	{
		if str.Value != nil {
			value = *str.Value
		}
	}

	if str.Optional && str.Value == nil {
		return str
	}

	if !reg.MatchString(value) {
		str.Add(str.Key, "mask", message, mask)
	}

	return str
}

func (str *String) ErrorMessage() Message {
	return str.Message
}

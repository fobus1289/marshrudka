package dependency

import (
	"reflect"
	"regexp"
	"strings"
	"test/util"
)

type IDependency interface {
	Set(service any) IDependency
	Get(service any) IDependency
	SetAll(services ...any) IDependency
	GetAll(services ...any) IDependency
	Fill(service any) IDependency
	FillAll(services ...any) IDependency
}

type IDependencyContainer interface {
	GetServices() map[reflect.Type]reflect.Value
	TrySearch(t reflect.Type, v *reflect.Value) bool
}

type dependency struct {
	Services map[reflect.Type]reflect.Value
}

func NewDependency() IDependency {
	return &dependency{
		Services: make(map[reflect.Type]reflect.Value),
	}
}

func (d *dependency) Set(s any) IDependency {

	value := util.Elem(s)
	{
		if value.Kind() == reflect.Invalid || !value.CanSet() || !value.CanAddr() {
			return d
		}
	}

	d.Services[value.Type()] = value

	if value.Kind() == reflect.Interface {

		interfaceValue := value.Elem()
		interfaceKind := interfaceValue.Kind()

		if interfaceKind == reflect.Invalid {
			return d
		}
		d.Services[interfaceValue.Type()] = interfaceValue

		if interfaceKind == reflect.Ptr {
			if self := interfaceValue.Elem(); self.Kind() != reflect.Invalid {
				d.Services[self.Type()] = self
			}
		}

	} else {
		addr := value.Addr()
		d.Services[addr.Type()] = addr
	}

	return d
}

func (d *dependency) Get(s any) IDependency {

	value := util.Elem(s)
	{
		if value.Kind() == reflect.Invalid || !value.CanSet() || !value.CanAddr() {
			return d
		}
	}

	if hasValue := d.Services[value.Type()]; hasValue.Kind() != reflect.Invalid {
		value.Set(hasValue)
	}

	return d
}

func (d *dependency) SetAll(services ...any) IDependency {
	for _, service := range services {
		d.Set(service)
	}
	return d
}

func (d *dependency) GetAll(services ...any) IDependency {

	for _, service := range services {
		d.Get(service)
	}

	return d
}

func (d *dependency) Fill(s any) IDependency {

	value := util.Elem(s)
	{
		if value.Kind() == reflect.Invalid || !value.CanSet() || !value.CanAddr() {
			return d
		}
	}

	if value.Kind() == reflect.Interface {
		value = util.Elem(value.Interface())
		{
			if value.Kind() == reflect.Invalid {
				return d
			}
		}
	}

	for i := 0; i < value.NumField(); i++ {

		field := value.Field(i)

		if val := d.Services[field.Type()]; val.Kind() != reflect.Invalid {
			field.Set(val)
		}
	}

	return d
}

func (d *dependency) FillAll(services ...any) IDependency {
	for _, service := range services {
		d.Fill(service)
	}
	return d
}

func (d *dependency) GetServices() map[reflect.Type]reflect.Value {
	return d.Services
}

func (d *dependency) TrySearch(t reflect.Type, v *reflect.Value) bool {

	if v == nil {
		return false
	}

	if hasValue := d.Services[t]; hasValue.Kind() != reflect.Invalid {
		*v = hasValue
		return true
	}

	return false
}

func (d *dependency) ParseFunc() {

}

var urlPath = "api/users/1/user1/edit"

func Split(i int) bool {
	paths := strings.Split(urlPath, "/")

	l := len(paths)

	if l < i {
		return false
	}

	_ = paths

	return true
}

func Count(i int) bool {
	count := strings.Count(urlPath, "/")

	if count < i {
		return false
	}

	_ = count

	return count > 0
}

var regExp = regexp.MustCompile(`^(api/users/(\w+)/(\w+)/edit)$`)

func Regexp(url string) bool {
	return regExp.MatchString(url)
}

func RegexpFindAll(url string) []string {
	return regExp.FindAllString(url, -1)
}

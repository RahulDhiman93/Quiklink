package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
)

// Form it creates a custom form struct and embeds a Url.Value object
type Form struct {
	url.Values
	Errors errors
}

// Valid Returns True for errors in form
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a new form
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required To be Used to validate the form fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field can not be blank")
		}
	}
}

// Has checked if forms field is in post or not
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		return false
	}
	return true
}

// MinLength is to be used to check minimum length of field
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be atleast %d character long", length))
		return false
	}
	return true
}

// IsEmail To use to check if email is in correct format
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}

// IsPhone To use to check if phone is in correct format
func (f *Form) IsPhone(field string) {
	if !govalidator.IsNumeric(f.Get(field)) || len(f.Get(field)) < 10 {
		f.Errors.Add(field, "Invalid phone number")
	}
}

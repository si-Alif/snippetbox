package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	NonFieldErrors []string
	FieldErrors map[string]string
}

// isValid function
func (v *Validator) Valid() bool{
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// add all the non field errors in the array
func (v *Validator) AddNonFieldError(message string){
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}


// add field errors to the map
func (v *Validator) AddFieldError(key string, message string){
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _ , exists := v.FieldErrors[key] ; !exists {
		v.FieldErrors[key] = message
	}
}

// if a field value is incorrect then add it to the map
func (v *Validator) CheckField(ok bool , key string , message string){
	if !ok {
		v.AddFieldError(key , message)
	}
}

func NotBlank(value string) bool{
	return strings.TrimSpace(value) != ""
}

func MaxChars(cnt int , s string) bool {
	return utf8.RuneCountInString(s) <= cnt
}

func MinChars (cnt int , s string) bool {
	return utf8.RuneCountInString(s) >= cnt
}


//comparing generic types and their counts
func PermittedValue[Comp comparable](vals Comp , permittedVals ...Comp) bool {
	return slices.Contains(permittedVals , vals)
}


// check if the provided string is a valid email address or not by comparing it with a regular expression
func Matches (value string , r *regexp.Regexp) bool{
	return r.MatchString(value)
}

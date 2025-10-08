package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// isValid function
func (v *Validator) Valid() bool{
	return len(v.FieldErrors) == 0
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


//comparing generic types and their counts 
func PermittedValue[Comp comparable](vals Comp , permittedVals ...Comp) bool {
	return slices.Contains(permittedVals , vals)
}

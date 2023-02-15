package api

import (
	"IMChat/utils"

	"github.com/go-playground/validator/v10"
)

var validGender validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if gender, ok := fieldLevel.Field().Interface().(int16); ok {
		return utils.IsSupportedGender(gender)
	}
	return false
}

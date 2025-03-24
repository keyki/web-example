package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"web-example/types"
	"web-example/util"
)

func Validate(val any) error {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := getFieldName(v.Type().Field(i))
		tags := v.Type().Field(i).Tag.Get("validate")
		if tags == "" {
			continue
		}
		for _, tag := range strings.Split(tags, ";") {
			if err := applyRule(tag, fieldName, field); err != nil {
				return err
			}
		}
	}
	return nil
}

func applyRule(tag, fieldName string, fieldValue reflect.Value) error {
	switch {
	case strings.HasPrefix(tag, "min="):
		num := getLengthValue(tag, "min=")
		if len(fieldValue.String()) < num {
			return fmt.Errorf("%s must be at least %d characters long", fieldName, num)
		}
	case strings.HasPrefix(tag, "max="):
		num := getLengthValue(tag, "max=")
		if len(fieldValue.String()) > num {
			return fmt.Errorf("%s must be max %d characters long", fieldName, num)
		}
	case strings.HasPrefix(tag, "required"):
		if len(fieldValue.String()) < 1 {
			return fmt.Errorf("%s is required", fieldName)
		}
	case strings.HasPrefix(tag, "oneOfRole="):
		if len(fieldValue.String()) > 0 {
			if !util.IsValidRole(fieldValue.String()) {
				return fmt.Errorf("%s must be one of %s", fieldName, types.Roles)
			}

		}
	}
	return nil
}

// min=2,max=10
func getLengthValue(tag, prefix string) int {
	num, _ := strconv.Atoi(strings.TrimPrefix(tag, prefix))
	return num
}

func getFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return field.Name
	}
	return jsonTag
}

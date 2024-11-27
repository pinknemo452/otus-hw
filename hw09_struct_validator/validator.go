package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type validationRule struct {
	Name  string
	Value string
}

type ValidationErrors []ValidationError

var (
	ErrNotAStructure            = errors.New("provided object is not a structure")
	ErrInvalidTag               = errors.New("invalid validate tag")
	ErrUnknownValidationKeyword = errors.New("unknown validation keyword")
	ErrLengthDoesNotMatch       = errors.New("string length does not match")
	ErrRegexpDoesNotMatch       = errors.New("string does not match regexp")
	ErrLessThanMinimum          = errors.New("int is less than minimum")
	ErrGreaterThanMaximum       = errors.New("int is greater than maximum")
	ErrNotInRange               = errors.New("field not in range")
)

const validateTag = "validate"

func (v ValidationErrors) Error() string {
	var err error
	for _, ve := range v {
		err = errors.Join(err, ve.Err)
	}

	return err.Error()
}

func Validate(v interface{}) error {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return ErrNotAStructure
	}
	val := reflect.ValueOf(v)
	validationErr := []error{}
	for i := range t.NumField() {
		value := val.Field(i)
		field := t.Field(i)
		tagValue, ok := field.Tag.Lookup(validateTag)
		if !ok {
			continue
		}
		err := processTag(value, tagValue)
		if err != nil {
			validationErr = append(validationErr, err)
		}
	}
	fmt.Println(validationErr)
	return errors.Join(validationErr...)
}

func processTag(fieldValue reflect.Value, tagValue string) error {
	rules, err := extractRules(tagValue)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		var err error
		switch rule.Name {
		case "len":
			err = validateLen(rule, fieldValue)
		case "regexp":
			err = validateRegexp(rule, fieldValue)
		case "in":
			err = validateIn(rule, fieldValue)
		case "min":
			err = validateMin(rule, fieldValue)
		case "max":
			err = validateMax(rule, fieldValue)
		default:
			err = ErrUnknownValidationKeyword
		}
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func extractRules(tagValue string) ([]validationRule, error) {
	rules := []validationRule{}
	if ands := strings.Split(tagValue, "|"); len(ands) > 1 {
		for _, and := range ands {
			if and == "" {
				return nil, ErrInvalidTag
			}
			rule, err := extractRule(and)
			if err != nil {
				return nil, err
			}
			rules = append(rules, rule)
		}
	} else {
		rule, err := extractRule(tagValue)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func extractRule(tagValue string) (validationRule, error) {
	if validators := strings.Split(tagValue, ":"); len(validators) == 2 {
		return validationRule{validators[0], validators[1]}, nil
	}
	return validationRule{}, ErrInvalidTag
}

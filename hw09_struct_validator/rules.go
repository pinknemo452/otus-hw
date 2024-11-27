package hw09structvalidator

import (
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func validateLen(rule validationRule, value reflect.Value) error {
	length, err := strconv.Atoi(rule.Value)
	if err != nil {
		return err
	}

	if value.Kind() == reflect.String {
		if len(value.String()) != length {
			return ErrLengthDoesNotMatch
		}
	} else if value.Kind() == reflect.Slice && value.Type().Elem().Kind() == reflect.String {
		for i := range value.Len() {
			if len(value.Index(i).String()) != length {
				return ErrLengthDoesNotMatch
			}
		}
	}

	return nil
}

func validateRegexp(rule validationRule, value reflect.Value) error {
	reg, err := regexp.Compile(rule.Value)
	if err != nil {
		return err
	}
	if value.Kind() == reflect.String {
		if !reg.Match([]byte(value.String())) {
			return ErrRegexpDoesNotMatch
		}
	} else if value.Kind() == reflect.Slice && value.Type().Elem().Kind() == reflect.String {
		for i := range value.Len() {
			if !reg.Match([]byte(value.Index(i).String())) {
				return ErrRegexpDoesNotMatch
			}
		}
	}
	return nil
}

func validateMax(rule validationRule, value reflect.Value) error {
	max, err := strconv.Atoi(rule.Value)
	if err != nil {
		return err
	}

	if value.Kind() == reflect.Int {
		if value.Int() > int64(max) {
			return ErrGreaterThanMaximum
		}
	} else if value.Kind() == reflect.Slice && value.Type().Elem().Kind() == reflect.Int {
		for i := range value.Len() {
			if value.Index(i).Int() > int64(max) {
				return ErrGreaterThanMaximum
			}
		}
	}

	return nil
}

func validateMin(rule validationRule, value reflect.Value) error {
	min, err := strconv.Atoi(rule.Value)
	if err != nil {
		return err
	}

	if value.Kind() == reflect.Int {
		if value.Int() < int64(min) {
			return ErrLessThanMinimum
		}
	} else if value.Kind() == reflect.Slice && value.Type().Elem().Kind() == reflect.Int {
		for i := range value.Len() {
			if value.Index(i).Int() < int64(min) {
				return ErrLessThanMinimum
			}
		}
	}

	return nil
}

func validateIn(rule validationRule, value reflect.Value) error {
	ruleSplitted := strings.Split(rule.Value, ",")
	if value.Kind() == reflect.String {
		if !slices.Contains(ruleSplitted, value.String()) {
			return ErrNotInRange
		}
		return nil
	} else if value.Kind() == reflect.Slice && value.Type().Elem().Kind() == reflect.String {
		for i := range value.Len() {
			if !slices.Contains(ruleSplitted, value.Index(i).String()) {
				return ErrNotInRange
			}
		}
		return nil
	}
	if len(ruleSplitted) != 2 {
		return ErrInvalidTag
	}
	lower, err := strconv.Atoi(ruleSplitted[0])
	if err != nil {
		return err
	}
	upper, err := strconv.Atoi(ruleSplitted[1])
	if err != nil {
		return err
	}
	if value.Kind() == reflect.Int {
		if !(value.Int() < int64(upper) && value.Int() > int64(lower)) {
			return ErrNotInRange
		}
	} else if value.Kind() == reflect.Slice && value.Type().Elem().Kind() == reflect.Int {
		for i := range value.Len() {
			if !(value.Index(i).Int() < int64(upper) && value.Index(i).Int() > int64(lower)) {
				return ErrNotInRange
			}
		}
	}
	return nil
}

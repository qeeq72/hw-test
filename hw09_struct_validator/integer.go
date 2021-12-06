package hw09structvalidator

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrValueIsUnderRange = errors.New("value is under range")
	ErrValueIsAboveRange = errors.New("value is above range")
	ErrValueIsNotInList  = errors.New("value is not in a list")
)

type IntegerValidator interface {
	Validate(int64) error
}

type IntegerMinValidator struct {
	Min int64
}

func (v IntegerMinValidator) Validate(i int64) error {
	if i < v.Min {
		return ErrValueIsUnderRange
	}
	return nil
}

type IntegerMaxValidator struct {
	Max int64
}

func (v IntegerMaxValidator) Validate(i int64) error {
	if i > v.Max {
		return ErrValueIsAboveRange
	}
	return nil
}

type IntegerListValidator struct {
	List []int64
}

func (v IntegerListValidator) Validate(i int64) error {
	for j := range v.List {
		if i == v.List[j] {
			return nil
		}
	}
	return ErrValueIsNotInList
}

type IntegerValidators []IntegerValidator

func (v IntegerValidators) Validate(i int64) []error {
	var errs []error
	for j := range v {
		err := v[j].Validate(i)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func InitIntegerValidator(validatorsMap map[string]string) (*IntegerValidators, error) {
	var validators IntegerValidators

	if listValuesStr, ok := validatorsMap["in"]; ok {
		listValuesArr := strings.Split(listValuesStr, ",")
		list := make([]int64, len(listValuesArr))
		for i := range listValuesArr {
			val, err := strconv.Atoi(listValuesArr[i])
			if err != nil {
				return nil, ErrInvalidValidatorTag
			}
			list[i] = int64(val)
		}
		validators = append(validators, IntegerListValidator{
			List: list,
		})
	}

	var min, max int64
	var isMin, isMax bool

	if minValueStr, ok := validatorsMap["min"]; ok {
		val, err := strconv.Atoi(minValueStr)
		if err != nil {
			return nil, ErrInvalidValidatorTag
		}
		min = int64(val)
		validators = append(validators, IntegerMinValidator{
			Min: min,
		})
		isMin = true
	}

	if maxValueStr, ok := validatorsMap["max"]; ok {
		val, err := strconv.Atoi(maxValueStr)
		if err != nil {
			return nil, ErrInvalidValidatorTag
		}
		max = int64(val)
		validators = append(validators, IntegerMaxValidator{
			Max: max,
		})
		isMax = true
	}

	if isMin && isMax {
		if min > max {
			return nil, ErrInvalidValidatorTag
		}
		if min == max {
			validators = []IntegerValidator{IntegerListValidator{
				List: []int64{min},
			}}
			return &validators, nil
		}
	}

	return &validators, nil
}

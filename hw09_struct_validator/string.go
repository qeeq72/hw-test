package hw09structvalidator

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrStringIsNotInList      = errors.New("string is not in a list")
	ErrStringHasWrongLength   = errors.New("string has wrong length")
	ErrStringIsUnmatchedRegex = errors.New("string is unmatched regexp")
)

type StringValidator interface {
	Validate(string) error
}

type StringLengthValidator struct {
	IsByteMode bool // считать длину в байтах, в противном случае в рунах
	Length     int
}

func (v StringLengthValidator) Validate(s string) error {
	if v.IsByteMode {
		if len(s) != v.Length {
			return ErrStringHasWrongLength
		}
	}
	if utf8.RuneCountInString(s) != v.Length {
		return ErrStringHasWrongLength
	}
	return nil
}

type StringRegexpValidator struct {
	Regexp *regexp.Regexp
}

func (v StringRegexpValidator) Validate(s string) error {
	if v.Regexp.MatchString(s) {
		return nil
	}
	return ErrStringIsUnmatchedRegex
}

type StringListValidator struct {
	List []string
}

func (v StringListValidator) Validate(s string) error {
	for j := range v.List {
		if s == v.List[j] {
			return nil
		}
	}
	return ErrStringIsNotInList
}

type StringValidators []StringValidator

func (v StringValidators) Validate(s string) []error {
	var errs []error
	for j := range v {
		err := v[j].Validate(s)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func InitStringValidator(validatorsMap map[string]string) (*StringValidators, error) {
	var validators StringValidators

	if strLengthStr, ok := validatorsMap["len"]; ok {
		val, err := strconv.Atoi(strLengthStr)
		if err != nil {
			return nil, ErrInvalidValidatorTag
		}
		validators = append(validators, StringLengthValidator{
			IsByteMode: false,
			Length:     val,
		})
	}

	if strLengthStr, ok := validatorsMap["blen"]; ok {
		val, err := strconv.Atoi(strLengthStr)
		if err != nil {
			return nil, ErrInvalidValidatorTag
		}
		validators = append(validators, StringLengthValidator{
			IsByteMode: true,
			Length:     val,
		})
	}

	if listValuesStr, ok := validatorsMap["in"]; ok {
		validators = append(validators, StringListValidator{
			List: strings.Split(listValuesStr, ","),
		})
	}

	if regexStr, ok := validatorsMap["regexp"]; ok {
		r, err := regexp.Compile(regexStr)
		if err != nil {
			return nil, ErrInvalidValidatorTag
		}
		validators = append(validators, StringRegexpValidator{
			Regexp: r,
		})
	}

	return &validators, nil
}

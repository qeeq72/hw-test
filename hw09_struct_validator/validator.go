package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrInvalidValidatorTag = errors.New("invalid validator tag")
	ErrNilPointer          = errors.New("nil pointer")
	ErrNotStructureValue   = errors.New("not a structure value")
	ErrUnsupportedType     = errors.New("unsupported type")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	n := len(v)
	if n < 1 {
		return ""
	}
	s := "Structure validation fail:\n"
	for i := 0; i < n; i++ {
		s += fmt.Sprintf("%s - %s;\n", v[i].Field, v[i].Err)
	}
	return s
}

func (v *ValidationErrors) Add(field string, err ...error) {
	for i := range err {
		*v = append(*v, ValidationError{
			Field: field,
			Err:   err[i],
		})
	}
}

func getValidatorMap(s string) (map[string]string, error) {
	validatorsStr := strings.Split(s, "|")

	validatorsMap := make(map[string]string)

	for j := range validatorsStr {
		sep := strings.Split(validatorsStr[j], ":")
		if len(sep) != 2 {
			return nil, ErrInvalidValidatorTag
		}
		if sep[1] == "" {
			continue
		}
		validatorsMap[sep[0]] = sep[1]
	}
	return validatorsMap, nil
}

func prevalidation(v interface{}) (reflect.Value, reflect.Type, error) {
	// 1. Получаем значение и тип
	value := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// 2. Проверяем на указатель
	if value.Kind() == reflect.Ptr {
		// 2.1. Если пустой, то просто выходим
		if value.IsNil() {
			return reflect.Value{}, nil, ErrNilPointer
		}

		// 2.2. Если не пустой, то берем у него значение и тип
		typ = value.Elem().Type()
		value = value.Elem()
	}

	// 3. Если итого получилась не структура, то выходим
	if value.Kind() != reflect.Struct {
		return reflect.Value{}, nil, ErrNotStructureValue
	}

	return value, typ, nil
}

func Validate(v interface{}) error { //nolint:gocognit
	// 1..3. Проверяем входные данные перед валидацией
	value, typ, err := prevalidation(v)
	if err != nil {
		return err
	}

	var validationErrs ValidationErrors

	// 4. Бежим по всем полям и валидируем их
	for i := 0; i < typ.NumField(); i++ {
		// 4.1. Берем значение и тип поля
		fieldValue := value.Field(i)
		fieldType := typ.Field(i)

		// 4.2. Проверяем наличие тега валидации, если нет - идем к следующему полю
		validationTag, ok := fieldType.Tag.Lookup("validate")
		if !ok {
			continue
		}

		// 4.3. Если тег есть, то берем из значения тега все, что считается "валидатором"
		validatorsMap, err := getValidatorMap(validationTag)
		if err != nil {
			return fmt.Errorf("%w: %v", err, fieldType.Name)
		}

		// 4.4. Проверяем поддерживаемые типы и валидимруем
		switch fieldValue.Kind() { //nolint:exhaustive
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			validators, err := InitIntegerValidator(validatorsMap)
			if err != nil {
				return fmt.Errorf("%w: %v", err, fieldType.Name)
			}
			validationErrs.Add(fieldType.Name, validators.Validate(fieldValue.Int())...)
		case reflect.String:
			validators, err := InitStringValidator(validatorsMap)
			if err != nil {
				return fmt.Errorf("%w: %v", err, fieldType.Name)
			}
			validationErrs.Add(fieldType.Name, validators.Validate(fieldValue.String())...)
		case reflect.Slice:
			length := fieldValue.Len()
			if length < 1 {
				continue
			}
			switch fieldValue.Index(0).Kind() { //nolint:exhaustive
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				validators, err := InitIntegerValidator(validatorsMap)
				if err != nil {
					return fmt.Errorf("%w: %v", err, fieldType.Name)
				}
				for i := 0; i < length; i++ {
					fieldValueElement := fieldValue.Index(i)
					validationErrs.Add(fmt.Sprintf("%v[%v]", fieldType.Name, i), validators.Validate(fieldValueElement.Int())...)
				}
			case reflect.String:
				validators, err := InitStringValidator(validatorsMap)
				if err != nil {
					return fmt.Errorf("%w: %v", err, fieldType.Name)
				}
				for i := 0; i < length; i++ {
					fieldValueElement := fieldValue.Index(i)
					validationErrs.Add(fmt.Sprintf("%v[%v]", fieldType.Name, i), validators.Validate(fieldValueElement.String())...)
				}
			default:
				validationErrs.Add(fieldType.Name, ErrUnsupportedType)
			}
		default:
			validationErrs.Add(fieldType.Name, ErrUnsupportedType)
		}
	}

	return validationErrs
}

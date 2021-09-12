package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func IsEmptyString(s string) bool {
	return len(s) == 0
}

func IsSlash(r rune) bool {
	return r == '\\'
}

func Unpack(s string) (string, error) {
	// 1.Проверяем строку на пустоту
	if IsEmptyString(s) {
		return "", nil
	}

	// 2. Заводим необходимые переменные
	var b strings.Builder // буфер для создания результирующей строки
	var letter rune       // для запоминания символа с предыдущей итерации
	var isShield bool     // для запоминания того, что символ текущей итерации экранирован

	// 3. Итерируемся по строке
	for _, r := range s {
		// 3.1. Если \, то проверяем экран это или просто символ
		if r == '\\' {
			if isShield {
				letter = r
				isShield = false
			} else {
				if letter != 0 {
					b.WriteRune(letter)
				}
				isShield = true
			}
			continue
		}

		// 3.2. Если цифра, то проверяем символ это или множитель для предыдущего символа
		var count int
		if unicode.IsDigit(r) {
			if isShield {
				letter = r
				isShield = false
			} else {
				if letter == 0 {
					return "", ErrInvalidString
				}
				count, _ = strconv.Atoi(string(r))
				b.WriteString(strings.Repeat(string(letter), count))
				letter = 0
			}
			continue
		}

		// 3.3. Если текущий символ экранированный, то он не цифра и не \, значит ошибка
		if isShield {
			return "", ErrInvalidString
		}

		// 3.4. Если символ не нулевой, то пишем его в буфер
		if letter != 0 {
			b.WriteRune(letter)
		}

		// 3.5. Запоминаем символ для последующей итерации
		letter = r
	}

	// 4. Проверяем и дописываем последний символ
	if isShield {
		return "", ErrInvalidString
	}
	if letter != 0 {
		b.WriteRune(letter)
	}

	return b.String(), nil
}

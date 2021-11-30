package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// 1. Просто читаем забираем список файлов из директории
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	// 2. Пробегаемся по файлам
	for _, file := range files {
		name := file.Name()

		// 2.1. Удаляем файлы с невалидными именами
		if strings.ContainsRune(name, '=') {
			continue
		}

		size := file.Size()

		// 2.2. Проверяем файл на пустоту
		if size == 0 {
			// 2.2.1. Файл пустой - помечаем переменную окружения на удаление
			env[name] = EnvValue{Value: "", NeedRemove: true}
		} else {
			// 2.2.2. Файл с данными - вычитываем данные и приводим их к требуемому формату
			data, err := os.OpenFile(dir+"/"+name, os.O_RDONLY, 0o777)
			if err != nil {
				return nil, err
			}
			b := make([]byte, size)
			_, err = data.Read(b)
			if err != nil {
				return nil, err
			}
			fileValue := strings.Split(string(b), "\n")[0]
			fileValue = strings.TrimRight(fileValue, " ")
			fileValue = strings.TrimRight(fileValue, "\t")
			fileValue = string(bytes.ReplaceAll([]byte(fileValue), []byte{0x00}, []byte{'\n'}))

			// 2.2.3. Пишем переменную в мапу
			env[name] = EnvValue{Value: fileValue, NeedRemove: false}
			data.Close()
		}
	}

	return env, nil
}

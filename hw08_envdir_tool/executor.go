package main

import (
	"errors"
	"os"
	"os/exec"
)

const ErrorCode = 1

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// 1. Сначала проверяем команду и аргументы
	if len(cmd) == 0 {
		// 1.1. Если команды нет, то и запускать нечего
		returnCode = ErrorCode
		return
	}
	// 1.2. Проводим все манипуляции с переменными окружения
	var err error
	for name := range env {
		err = os.Unsetenv(name)
		if err != nil {
			returnCode = ErrorCode
			return
		}
		if !env[name].NeedRemove {
			err = os.Setenv(name, env[name].Value)
			if err != nil {
				returnCode = ErrorCode
				return
			}
		}
	}

	command := cmd[0]
	app := exec.Command(command)
	app.Env = os.Environ()

	// 1.3. Пробрасываем потоки ввода, вывода и ошибок
	app.Stdin = os.Stdin
	app.Stdout = os.Stdout
	app.Stderr = os.Stderr

	// 1.4. Добавляем аргументы, если они есть
	if len(cmd) > 1 {
		app.Args = append(app.Args, cmd[1:]...)
	}

	// 1.5. Выполняем программу
	if err := app.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			returnCode = exitErr.ExitCode()
		}
	}
	return
}

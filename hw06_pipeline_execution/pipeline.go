package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// 0. Проверяем входные данные на nil
	if in == nil {
		return nil
	}
	if stages == nil {
		return in
	}

	// 1. Даем возможность прервать до начала выполнения stages
	out := stageWithInterruption(done, in)

	// 2. Запускаем выполнение stages
	for _, stage := range stages {
		// 2.1. Если попался пустой stage, то пропускаем
		if stage == nil {
			continue
		}
		// 2.2. В противном случае выполняем
		out = stage(stageWithInterruption(done, out))
	}
	return out
}

// Создаем функцию для прерывания работы.
func stageWithInterruption(done In, in In) Out {
	// 0. Создаем канал, в который будем писать вход и передавать дальше на чтение
	out := make(Bi)

	// 1. В отдельной горутине проверяем done, либо просто передаем канал дальше
	go func() {
		defer close(out)
		for {
			// 1.1 Создаем приоритет для done
			select {
			case <-done:
			default:
			}

			// 1.2 А тут ловим данные со входа или отрабатываем done
			select {
			case <-done:
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				out <- value
			}
		}
	}()

	return out
}

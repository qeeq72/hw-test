package hw03frequencyanalysis

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type wordInText struct {
	Word  string
	Count int
}

var (
	reRemovePrefix = regexp.MustCompile(`^[^a-zA-Zа-яА-Я0-9]+`)
	reRemoveSuffix = regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9]+$`)
)

func Top10(text string) []string {
	// 1. Проверяем строку на пустоту
	if text == "" {
		fmt.Println("Empty string")
		return nil
	}

	// 2. Разбиваем текст на отдельные слова и сразу проверяем количество слов
	singleWords := strings.Fields(text)
	if len(singleWords) < 10 {
		fmt.Println("Less than 10 words")
		return nil
	}

	// 3. Создаем словарь и заполняем его уникальным словами с их повторениями
	m := map[string]int{}
	for _, word := range singleWords {
		// 3.1. Проверяем на пустоту и тире
		if word == "-" || word == "" {
			continue
		}

		// 3.2. Отсекаем знаки препинания до и после слова
		word = reRemovePrefix.ReplaceAllString(word, "")
		word = reRemoveSuffix.ReplaceAllString(word, "")

		// 3.3. Переводим буквы в нижний регистр и пишем слово в мапу
		word = strings.ToLower(word)
		m[word]++
	}

	// 4. Проверяем заполненность мапы
	if len(m) < 10 {
		fmt.Println("Less than 10 unique words")
		return nil
	}

	// 5. Перекладываем мапу в слайс структур
	wordsInText := make([]wordInText, 0, 10)
	for k, v := range m {
		wordsInText = append(wordsInText, wordInText{Word: k, Count: v})
	}

	// 6. Сортируем слайс структур по количеству повторения и лексикографически
	sort.Slice(wordsInText, func(i, j int) bool {
		if wordsInText[i].Count > wordsInText[j].Count {
			return true
		}
		if wordsInText[i].Count < wordsInText[j].Count {
			return false
		}
		ri := []rune(wordsInText[i].Word)
		rj := []rune(wordsInText[j].Word)
		rMinLen := len(ri)
		if len(rj) < rMinLen {
			rMinLen = len(rj)
		}
		if rMinLen == 0 {
			return false
		}
		for k := 0; k < rMinLen; k++ {
			if ri[k] < rj[k] {
				return true
			}
			if ri[k] > rj[k] {
				return false
			}
		}
		return false
	})

	// 7. Создаем и заполняем итоговый слайс
	topWords := make([]string, 0, 10)

	for i := 0; i < 10; i++ {
		topWords = append(topWords, wordsInText[i].Word)
	}

	return topWords
}

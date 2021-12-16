package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

/*
	Оптимизации getUsers:
	1. Использована другая библиотека для ускорения анмаршалинга JSON
	2. Для экономии памяти используется Scanner
	3. Не используется каст строки к слайсу байт на входе в функцию анмаршалинга
*/
func getUsers(r io.Reader) (result users, err error) {
	br := bufio.NewScanner(r)
	// json := jsoniter.ConfigCompatibleWithStandardLibrary
	json := jsoniter.ConfigFastest

	var count int
	for br.Scan() {
		var user User
		if err = json.Unmarshal(br.Bytes(), &user); err != nil {
			return
		}
		result[count] = user
		count++
	}

	return
}

/*
	Оптимизации countDomains:
	1. Используется скомпилированная регулярка
	2. Добавились некоторые проверки на корректность домена
*/
func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	re, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	var num int
	for i := range u {
		matched := re.Match([]byte(u[i].Email))

		if matched {
			sep := strings.SplitN(u[i].Email, "@", 3)
			if len(sep) != 2 {
				continue
			}
			find := strings.ToLower(sep[1])
			sep = strings.SplitN(u[i].Email, ".", 3)
			if len(sep) != 2 {
				continue
			}
			num = result[find]
			num++
			result[find] = num
		}
	}
	return result, nil
}

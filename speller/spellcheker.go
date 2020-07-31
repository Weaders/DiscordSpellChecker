package speller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/thoas/go-funk"
)

type badWord struct {
	Word string   `json:word`
	S    []string `json:s`
}

// Check string on error
func CheckString(str string) []string {

	resp, err := http.Get("https://speller.yandex.net/services/spellservice.json/checkText?text=" + url.PathEscape(str))

	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		data := []badWord{}
		json.Unmarshal(body, &data)

		result := (funk.Map(data, func(word badWord) string {

			var result strings.Builder

			result.WriteString("Говно слово: " + word.Word + "\r\n")
			result.WriteString("Правильно:  " + strings.Join(word.S, ","))

			return result.String()

		})).([]string)

		return result

	}

	return []string{"Ошибка, Яндекс не болей"}

}

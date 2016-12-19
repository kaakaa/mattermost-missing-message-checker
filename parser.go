package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	zglob "github.com/mattn/go-zglob"
)

var expJSX = []*regexp.Regexp{
	regexp.MustCompile(`(?s)<FormattedMessage.*?id='(.*?)'.*?defaultMessage=(.*?)[\s]*?/>`),
	regexp.MustCompile(`(?s)<FormattedHTMLMessage.*?id='(.*?)'.*?defaultMessage=(.*?)[\s]*?/>`),
	regexp.MustCompile(`(?s)Utils.localizeMessage\('(.*?)', '(.*?)'\)`),
}

var expGo = []*regexp.Regexp{
	regexp.MustCompile(`utils.T\("(.*?)"`),
	regexp.MustCompile(`c.T\("(.*?)"`),
	regexp.MustCompile(`model.NewLocAppError\(".*?", "(.*?)"`),
}

func parseFrontI18N(path string) (map[string]string, error) {
	b, err := readFile(path)
	if err != nil {
		return nil, err
	}

	var ret map[string]string
	json.Unmarshal(b, &ret)
	return ret, nil
}

func parseServerI18N(path string) ([]map[string]interface{}, error) {
	b, err := readFile(path)
	if err != nil {
		return nil, err
	}

	var ret []map[string]interface{}
	json.Unmarshal(b, &ret)
	return ret, nil

}

func parseJSX() ([]Message, error) {
	matches, err := zglob.Glob(`./platform/webapp/**/*.jsx`)
	if err != nil {
		return nil, err
	}
	var messages []Message
	for _, v := range matches {
		fmt.Printf("Parsing %s\n", v)
		for _, e := range expJSX {
			for _, r := range []regexp.Regexp{*e} {
				m, err := parse(v, r, true)
				if err != nil {
					return nil, err
				}
				messages = append(messages, m...)
			}
		}
	}

	return messages, nil
}

func parseGo() ([]Message, error) {
	matches, err := zglob.Glob(`./platform/**/*.go`)
	if err != nil {
		return nil, err
	}
	var messages []Message
	for _, v := range matches {
		if strings.Contains(v, "/vendor/") {
			continue
		}

		fmt.Printf("Parsing %s\n", v)
		for _, e := range expGo {
			for _, r := range []regexp.Regexp{*e} {
				m, err := parse(v, r, false)
				if err != nil {
					return nil, err
				}
				messages = append(messages, m...)
			}
		}
	}
	return messages, nil
}

func parse(path string, exp regexp.Regexp, hasDefaultMesssage bool) ([]Message, error) {
	f, err := readFile(path)
	if err != nil {
		return nil, err
	}
	var messages []Message
	for _, b := range exp.FindAllSubmatch(f, -1) {
		m := Message{}
		m.ID = string(b[1])
		if hasDefaultMesssage {
			v := string(b[2])
			m.DefaultMessage = v[1 : len(v)-1]
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}

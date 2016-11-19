package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	zglob "github.com/mattn/go-zglob"
)

var exp1 = regexp.MustCompile(`(?s)<FormattedMessage.*?id='(.*?)'.*?defaultMessage=(.*?)[\s]*?/>`)

func parseI18N() (map[string]string, error) {
	path := "./platform/webapp/i18n/en.json"
	b, err := readFile(path)
	if err != nil {
		return nil, err
	}

	var ret map[string]string
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
		for _, r := range []regexp.Regexp{*exp1} {
			m, err := parse(v, r)
			if err != nil {
				return nil, err
			}
			messages = append(messages, m...)
		}
	}
	return messages, nil
}

func parse(path string, exp regexp.Regexp) ([]Message, error) {
	b, err := readFile(path)
	if err != nil {
		return nil, err
	}
	var messages []Message
	for _, b := range exp.FindAllSubmatch(b, -1) {
		v := string(b[2])
		messages = append(messages, Message{ID: string(b[1]), DefaultMessage: v[1 : len(v)-1]})
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

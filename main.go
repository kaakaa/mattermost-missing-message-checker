package main

import (
	"encoding/json"
	"os"
)

type Message struct {
	ID             string
	DefaultMessage string
}

type Diff struct {
	MissingIDinJson   []Message
	DifferentMessages []Message
}

func main() {
	jsonMessage, err := parseI18N()
	if err != nil {
		panic(err)
	}
	jsxMessage, err := parseJSX()
	if err != nil {
		panic(err)
	}

	var ret Diff
	for _, v := range jsxMessage {
		if val, ok := jsonMessage[v.ID]; ok {
			if v.DefaultMessage != val {
				ret.DifferentMessages = append(ret.DifferentMessages, v)
			}
		} else {
			ret.MissingIDinJson = append(ret.MissingIDinJson, v)
		}
	}

	for _, v := range ret.MissingIDinJson {
		jsonMessage[v.ID] = v.DefaultMessage
	}
	f, err := os.OpenFile("./new_en.json", os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err = enc.Encode(jsonMessage)
	if err != nil {
		panic(err)
	}
}

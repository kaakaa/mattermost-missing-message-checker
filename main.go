package main

import (
	"encoding/json"
	"log"
	"os"
)

var (
	serverMessageFilePath = "./mattermost-server/i18n/en.json"
	webappMessageFilePath = "./mattermost-webapp/i18n/en.json"
	mobileMessageFilePath = "./mattermost-mobile/assets/base/i18n/en.json"
)

var (
	serverSourceFilePattern = "./mattermost-server/**/*.go"
	webappSourceFilePattern = "./mattermost-webapp/**/*.jsx"
	mobileSourceFilePattern = "./mattermost-mobile/**/*.js"
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
	parseFront()
	parseServer()
	parseMobile()
}

func parseFront() {
	jsonMessage, err := parseFrontI18N(webappMessageFilePath)
	if err != nil {
		panic(err)
	}
	jsxMessage, err := parseJSX(webappSourceFilePattern)
	log.Println("hogehoge2")
	log.Println(len(jsxMessage))
	if err != nil {
		panic(err)
	}
	jsxResult := mergeFront(jsonMessage, jsxMessage)
	write(jsxResult, "webi18n.json")
}

func parseMobile() {
	jsonMessage, err := parseFrontI18N(mobileMessageFilePath)
	if err != nil {
		panic(err)
	}
	jsxMessage, err := parseRN(mobileSourceFilePattern)
	log.Println("hogehoge")
	log.Println(len(jsxMessage))
	if err != nil {
		panic(err)
	}
	jsxResult := mergeFront(jsonMessage, jsxMessage)
	write(jsxResult, "mobile18n.json")
}

func parseServer() {
	jsonMessage, err := parseServerI18N(serverMessageFilePath)
	if err != nil {
		panic(err)
	}
	goMessage, err := parseGo(serverSourceFilePattern)
	if err != nil {
		panic(err)
	}
	goResult := mergeServer(jsonMessage, goMessage)
	write(goResult, "serveri18n.json")
}

func mergeFront(origin map[string]string, coding []Message) map[string]string {
	var ret Diff
	for _, v := range coding {
		if val, ok := origin[v.ID]; ok {
			if v.DefaultMessage != val {
				ret.DifferentMessages = append(ret.DifferentMessages, v)
			}
		} else {
			ret.MissingIDinJson = append(ret.MissingIDinJson, v)
		}
	}

	for _, v := range ret.MissingIDinJson {
		origin[v.ID] = v.DefaultMessage
	}
	return origin
}

func mergeServer(origin []map[string]interface{}, coding []Message) []map[string]interface{} {
	var originIDs []string
	for _, v := range origin {
		originIDs = append(originIDs, v["id"].(string))
	}
	var ret Diff
	for _, v := range coding {
		if !stringInSlice(v.ID, originIDs) {
			ret.MissingIDinJson = append(ret.MissingIDinJson, v)
		}
	}

	for _, v := range ret.MissingIDinJson {
		m := map[string]interface{}{
			"id":        v.ID,
			"translate": "",
		}
		origin = append(origin, m)
	}
	return origin
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func write(messages interface{}, path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err = enc.Encode(messages)
	if err != nil {
		return err
	}
	return nil
}

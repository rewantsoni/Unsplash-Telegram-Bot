package main

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"net/http"
	"os"
)


//Create the InlineKeyboard
func newInlineKeyboard(data data) tgbotapi.InlineKeyboardMarkup {
	MarshalData, _ := json.Marshal(data)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("More Like this", string(MarshalData)),
		),
	)
}

//New Request
func newPhotoRequest(d data) *http.Request {
	req, err := http.NewRequest("GET", "https://api.unsplash.com/photos/random", nil)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	q := req.URL.Query()
	q.Add("client_id", os.Getenv("CLIENT_ID"))
	q.Add("count", "1")
	if !d.Random {
		q.Add("query",d.Query)
	}
	req.URL.RawQuery = q.Encode()
	return req
}

//Fetch Response Object
func getResponse(req *http.Request) ([]response, error) {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []response{}, errors.New("request failed")
	}
	data, err := ioutil.ReadAll(res.Body)
	if string(data)=="Rate Limit Exceeded" {
		return []response{}, errors.New("rate limit exceeded please try after sometime")
	}
	if err != nil {
		return []response{}, errors.New("couldn't read response")
	}
	defer res.Body.Close()

	var parsedRes []response
	err = json.Unmarshal(data,&parsedRes)
	if err != nil {
		return []response{}, errors.New("couldn't find a image")
	}
	return parsedRes, nil
}

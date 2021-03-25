package main

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
)

//EV put new slack events
var EV *slack.MessageEvent

//RTM use for sending events to slack
var RTM *slack.RTM

//BotToken Put your slackbot token here
const BotToken string = "xoxb-3069876617-1921979534688-4VgW7DuSzbvFRzqnW29ljbdw"

//DefaultChannel Put your default channel
const DefaultChannel string = "#replica-memo"

func main() {
	var api *slack.Client = slack.New(BotToken)

	RTM = api.NewRTM()

	go RTM.ManageConnection()

	for msg := range RTM.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			fmt.Printf("Start connection with Slack\n")
		case *slack.MessageEvent:
			EV = ev
			ListenTo()
		}
	}
}

//ListenTo excute functions under suitable conditions
func ListenTo() {
	switch {
	case strings.Contains("こんにちは", EV.Text):
		RTM.SendMessage(RTM.NewOutgoingMessage("こんにちは。<@"+EV.User+">さん", EV.Channel))
		return
	}
}
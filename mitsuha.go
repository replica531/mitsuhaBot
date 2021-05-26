package main

import (
    "fmt"
	"github.com/slack-go/slack"
)

func main() {

	RTM = api.NewRTM()

	go RTM.ManageConnection()

    for{
        for msg := range RTM.IncomingEvents {
            switch ev := msg.Data.(type) {
                case *slack.ConnectedEvent:
                    fmt.Printf("Start connection with Slack\n")
                case *slack.MessageEvent:
                    EV = ev
                    ListenTo()
                default:
                    Reglarsend()
            }
        }
    }
}

//scp -r mitsuhaBot replica@kmc.gr.jp:~/Program
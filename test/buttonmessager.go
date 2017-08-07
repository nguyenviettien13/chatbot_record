package main

import (
	"github.com/michlabs/fbbot"
	"fmt"
	"log"
)

const (
	PAGEACCESS_TOKENn = "EAAEo7uCJDOEBAIM0NNK8LwfJ25YHOPkKqZCiVsCowsjMUdLUB2l0ABXTZCkZBIM5rFzzTmPvqsdxKmOBMl2P4ZAwxa5qe2Fgk1w6XF34SMOmvbzllwaN9HUIsObcxxBkikkp4ApNo0ceHOIgvhE25B3DiqBZBipgQskeDkBvOsZBWQTFuN8h6y"
	VERIFY_TOKENn    = "1234"
	PORTn            = 8080
)

type Recorder struct {}

func (r Recorder) HandleMessage(bot *fbbot.Bot, msg *fbbot.Message) {
	globalMap["name"] = "minhnd"
	fmt.Println("prepare sending ")
	b := fbbot.NewButtonMessage()
	b.Text ="Bạn có muốn tiếp tục?"
	b.AddPostbackButton("Có","yes")
	b.AddPostbackButton("Không","no")

	bot.Send(msg.Sender, b)
	fmt.Println("sent ")
}

func (r Recorder) HandlePostback(bot *fbbot.Bot, msg *fbbot.Postback) {
	switch msg.Payload {
	case "yes":
		bot.SendText(msg.Sender, "You choose " + msg.Payload)
	case "no":
		bot.SendText(msg.Sender, "Goodbye")
	default:
		log.Println("Switch case does not exist")
	}

}

type UserId string

type TimeNumber int
var globalMap map[string]string


func main() {
	var r Recorder
	globalMap = make(map[string]string)
	//init a bot
	bot := fbbot.New(PORTn, VERIFY_TOKENn, PAGEACCESS_TOKENn)

	//bot call AddMessageHandle
	bot.AddMessageHandler(r)
	bot.AddPostbackHandler(r)
	bot.Run()
}




package main

import (
	"flag"
	"fmt"

	"github.com/scristofari/sms-sender/twilio"
)

var (
	to      string
	content string
)

func main() {
	flag.StringVar(&to, "to", "", "Get the phone number")
	flag.StringVar(&content, "content", "", "Get the content of the SMS")
	flag.Parse()

	if err := twilio.SendSMS(to, content); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Done !!")
}

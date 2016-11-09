package twilio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	validate "gopkg.in/go-playground/validator.v9"
)

var (
	twilioPatternURL = "https://api.twilio.com/2010-04-01/Accounts/%s/Messages"
)

type sms struct {
	From    string `validate:"required,contains=+"`
	To      string `validate:"required,contains=+"`
	Content string `validate:"required"`
}

// SendSMS send the sms.
// This method is called by the cli or the server.
// To: +33XXXXXXXXX.
// Content: the message.
func SendSMS(to, content string) error {
	envs, err := getTwilioEnvs()
	if err != nil {
		return err
	}

	s := sms{
		From:    envs["TWILIO_SMS_FROM"],
		To:      to,
		Content: content,
	}

	validator := validate.New()
	err = validator.Struct(&s)
	if err != nil {
		return fmt.Errorf("failed to validate the SMS, %s", err)
	}

	data := new(bytes.Buffer)
	err = json.NewEncoder(data).Encode(s)
	if err != nil {
		return fmt.Errorf("failed to encode the SMS, %s", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(twilioPatternURL, envs["TWILIO_ACCOUNT_SID"]),
		data,
	)
	if err != nil {
		return fmt.Errorf("failed to create a new request : %s", err)
	}
	req.SetBasicAuth(envs["TWILIO_ACCOUNT_SID"], envs["TWILIO_AUTH_TOKEN"])

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send the request : %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get a valid response, get %d", res.StatusCode)
	}

	return nil
}

func getTwilioEnvs() (map[string]string, error) {
	envs := []string{
		"TWILIO_SMS_FROM",
		"TWILIO_ACCOUNT_SID",
		"TWILIO_AUTH_TOKEN",
	}

	m := map[string]string{}
	for _, key := range envs {
		v := os.Getenv(key)
		if v == "" {
			return nil, fmt.Errorf("environment variable %s is missing", key)
		}
		m[key] = v
	}
	return m, nil
}

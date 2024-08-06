package workers

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"urbangrid.com/functions"
)

func SendSMS(ctx context.Context, t *asynq.Task) error {
	accountSid := "AC9294c0f603f082aaffb31d6918d0cc6f"
	authToken := "35d12e0a8aa047661024fea954c58e21"

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo("+916353940369")
	params.SetFrom("+19388391544")
	params.SetBody(string(t.Payload()))

	resp, err := client.Api.CreateMessage(params)

	if err != nil {
		fmt.Println(err)
		return err
	} else {
		functions.CreateSMSLog(*resp.Body, *resp.Sid)
		return nil
	}
}

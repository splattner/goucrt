package client

import (
	"log"
	"time"

	"github.com/splattner/goucrt/pkg/entities"
	"github.com/splattner/goucrt/pkg/integration"
)

type Client struct {
	IntegrationDriver *integration.Integration
}

func NewClient(i *integration.Integration) *Client {

	client := Client{}

	client.IntegrationDriver = i

	return &client

}

func (c *Client) SetupClient() {

	c.IntegrationDriver.SetHandleSetupFunction(c.HandleSetup)

}

func (c *Client) HandleSetup() {

	var userAction = integration.RequiredUserAction{
		Confirmation: integration.ConfirmationPage{
			Title: integration.LanguageText{
				En: "You are about to add this integration",
			},
		},
	}

	// Start the setup
	c.IntegrationDriver.SetDriverSetupState(integration.StartEvent, integration.SetupState, integration.NoneError, &userAction)

	// For Testing, just Wait a bit
	time.Sleep(1 * time.Second)

	// Finish the setup
	c.IntegrationDriver.SetDriverSetupState(integration.StopEvent, integration.OkState, integration.NoneError, nil)

	// Some dummy test data
	button := entities.NewButtonEntity("mybutton", entities.LanguageText{En: "My Button", De: "Mein Button"}, "")
	button.AddCommand(entities.PushButtonEntityCommand, c.HandleButtonPressCommand)
	c.IntegrationDriver.AddEntity(button)

}

func (c *Client) HandleButtonPressCommand(button entities.ButtonEntity) {
	log.Println("Button " + button.Id + "pressed")
}

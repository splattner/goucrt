package integration

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/splattner/goucrt/pkg/integration/entities"
)

type integration struct {
	DeviceId string

	Metadata DriverMetadata

	authToken string

	DeviceState string

	listenPort int
	websocket  *websocket.Conn

	Remote remote

	Entities []entities.Entity

	// User input result of a SettingsPage as key values.
	// key: id of the field
	// value: entered user value as string. This is either the entered text or number, selected checkbox state or the selected dropdown item id.
	//⚠️ Non native string values as numbers or booleans are represented as string values!
	UserInputValues map[string]string

	SubscribedEntities []string
}

func NewIntegration() (*integration, error) {

	metadata := DriverMetadata{
		DriverId: "myintegraiton",
		Name: LanguageText{
			En: "My UCRT Integration",
		},
		Version: "0.0.1",
	}

	i := integration{
		listenPort: 8080,
		Metadata:   metadata,
	}

	return &i, nil

}

func (i *integration) Run() error {

	http.HandleFunc("/ws", i.wsEndpoint)

	listenAddress := fmt.Sprintf(":%d", i.listenPort)

	log.Fatal(http.ListenAndServe(listenAddress, nil))

	return nil

}

func (i *integration) addEntity(e entities.Entity) {
	i.Entities = append(i.Entities, e)
}

package graphql

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	graphql "github.com/graph-gophers/graphql-go"
)

type SocketHandler struct {
	Schema   *graphql.Schema
	Upgrader websocket.Upgrader
}

func (h *SocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	connection, upgradeError := h.Upgrader.Upgrade(w, r, nil)

	if upgradeError != nil {
		log.Println(upgradeError)
		return
	}

	for {
		msgType, message, error := connection.ReadMessage()

		if error != nil {
			return
		}

		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}

		if err := json.Unmarshal(message, &params); err != nil {
			errorResponse, _ := json.Marshal(err)
			connection.WriteMessage(msgType, errorResponse)
		}

		response := h.Schema.Exec(r.Context(), params.Query, params.OperationName, params.Variables)
		responseJSON, err := json.Marshal(response)
		if err != nil {
			errorResponse, _ := json.Marshal(err)
			connection.WriteMessage(msgType, errorResponse)
		}

		if error = connection.WriteMessage(msgType, responseJSON); error != nil {
			errorResponse, _ := json.Marshal(error)
			connection.WriteMessage(msgType, errorResponse)
		}
	}
}

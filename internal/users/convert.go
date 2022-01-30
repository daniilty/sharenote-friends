package users

import (
	"fmt"

	events "github.com/daniilty/sharenote-kafka-events"
)

func eventDataToUserDeleteEventData(data map[string]interface{}) (*events.UserDeleteEvent, error) {
	var ok bool
	u := &events.UserDeleteEvent{}

	u.ID, ok = data["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	return u, nil
}

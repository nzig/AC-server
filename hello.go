package hello

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"cloud.google.com/go/pubsub"
)

const projectName string = "kodicloud-169614"
const topicName string = "AirCon"

func init() {
	http.HandleFunc("/send", send)
}

func send(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if !isAllowedUser(u.Email) {
		errMsg := fmt.Sprintf("%s Not Allowed", u.Email)
		http.Error(w, errMsg, http.StatusForbidden)
		return
	}

	temperature := r.FormValue("temp")
	log.Infof(c, "[%s] executed %s", u.Email, temperature)

	topic, err := getTopic(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := topic.Publish(c, &pubsub.Message{
		Data: []byte(temperature),
		Attributes: map[string]string{
			"user": u.Email,
		},
	})
	_, err = result.Get(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Sent %s", template.HTMLEscapeString(temperature))
}

func isAllowedUser(email string) bool {
	switch email {
	case "nadavz0@gmail.com",
		"dor.israeli@gmail.com",
		"talflom@gmail.com":
		return true
	}
	return false
}

var gTopic *pubsub.Topic
var gClient *pubsub.Client

func getTopic(ctx context.Context) (*pubsub.Topic, error) {
	if gClient == nil {
		client, err := pubsub.NewClient(ctx, projectName)
		if err != nil {
			log.Errorf(ctx, "Error creating client: %v", err)
			return nil, errors.New("Failed to create pubsub client")
		}
		gClient = client
	}

	if gTopic == nil {
		topic := client.Topic(topicName)
		exists, err := topic.Exists(ctx)
		if err != nil {
			log.Errorf(ctx, "Error checking for topic: %v", err)
			return nil, errors.New("Failed checking for pubsub topic")
		}
		if !exists {
			newTopic, err := gClient.CreateTopic(ctx, topicName)
			if err != nil {
				log.Errorf(ctx, "Failed to create topic: %v", err)
				return nil, errors.New("Failed creating pubsub topic")
			}
			gTopic = newTopic
		}
		gTopic = topic
	}

	return gTopic, nil
}

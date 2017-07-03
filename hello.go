package hello

import (
	"fmt"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"cloud.google.com/go/pubsub"
)

const projectName string = "kodicloud-169614"

func init() {
	http.HandleFunc("/send", send)
}

func send(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	switch u.Email {
	case "nadavz0@gmail.com":
	case "dor.israeli@gmail.com":
	case "giljpeled@gmail.com":

	default:
		http.Error(w, "Not allowed", http.StatusForbidden)
		return
	}
	temp := r.FormValue("temp")
	log.Infof(c, "[%s] executed %s", u.Email, temp)
	client, err := pubsub.NewClient(c, projectName)
	if err != nil {
		http.Error(w, "Failed to create client", http.StatusForbidden)
		return
	}
	topic := client.Topic("AirCon")
	if topic == nil {
		http.Error(w, "Failed to get topic", http.StatusForbidden)
		return
	}
	result := topic.Publish(c, &pubsub.Message{Data: []byte(temp)})
	_, err = result.Get(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Sent %s", template.HTMLEscapeString(temp))
}

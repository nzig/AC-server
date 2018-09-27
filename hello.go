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
	if !isAllowedUser(u.Email) {
		errMsg := fmt.Sprintf("%s Not Allowed", u.Email)
		http.Error(w, errMsg, http.StatusForbidden)
		return
	}

	temperature := r.FormValue("temp")
	log.Infof(c, "[%s] executed %s", u.Email, temperature)
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

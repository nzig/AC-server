package hello

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"

	"cloud.google.com/go/pubsub"
)

const projectName string = "kodicloud-169614"

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/send", send)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, mainHTML)
}

const mainHTML = `
<html>
	<body>
		<form action="/send" method="post">
		<select name="temp">
			<option>off</option>
			<option>16</option>
			<option>17</option>
			<option>18</option>
			<option>19</option>
			<option>20</option>
			<option>21</option>
			<option>22</option>
			<option>23</option>
			<option>24</option>
			<option>25</option>
			<option>26</option>
			<option>27</option>
			<option>28</option>
			<option>29</option>
			<option>30</option>
		</select> 
		<input type="submit" value="Send">
		</form>
	</body>
</html>
`

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
	log.Printf("[%s] executed %s", u.Email, temp)
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

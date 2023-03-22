package webhandler

import (
	"encoding/json"
	"net/http"

	"github.com/philcantcode/localmapper/feeds"
)

type FeedHandler struct {
}

var Feed = FeedHandler{}

func (feed *FeedHandler) HTTP_JSON_GetWordlists(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(feeds.GetAllWordlists())
}

package scheduler

import (
	"awesomeProject/scheduler/dbops"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func vidDelRecHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")

	if len(vid) == 0 {
		SendResponse(w, 400, "video id should not be empty")
		return
	}

	err := dbops.AddVideoDeletionRecord(vid)
	if err != nil {
		SendResponse(w, 500, "Internal server error")
		return
	}

	SendResponse(w, 200, "")
	return
}

package streamserver

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	vl := VIDEO_DIR + vid

	video, err := os.Open(vl)
	if err != nil {
		log.Printf("Error when trying to open a file: %v", err)
		SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)

	defer video.Close()

}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "File is too big")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file error: %v", err)
		SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
	}

	filename := p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR+filename, data, 0666)
	if err != nil {
		log.Printf("ERror when writing to file: %v", err)
		SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "uploaded successfully")

}

// func testPageHandler(w http.ResponseWriter, r *http.Request) {
// 	t, _ := template.ParseFiles("./video/upload.html")
// 	t.Execute(w, nil)
// }

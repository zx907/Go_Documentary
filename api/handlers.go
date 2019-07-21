package main

import (
	"awesomeProject/api/dbops"
	"awesomeProject/api/defs"
	"awesomeProject/api/session"
	"awesomeProject/api/utils"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &defs.UserCredential{}

	if err := json.Unmarshal(res, ubody); err != nil {
		SendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		SendErrorResponse(w, defs.ErrorDBError)
		return
	}

	id := session.GenerateNewSessionId(ubody.Username)
	su := &defs.SignedUp{Success: true, SessionId: id}

	if resp, err := json.Marshal(su); err != nil {
		SendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		SendNormalResponse(w, string(resp), 201)
	}
}

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	log.Printf("%s", res)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		log.Printf("%s", err)
		SendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	// Validate the request body
	uname := p.ByName("username")
	log.Printf("Login url name: %s", uname)
	log.Printf("Login body name: %s", ubody)
	if uname != ubody.Username {
		SendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}

	log.Printf("%s", ubody.Username)
	pwd, err := dbops.GetUserCredential(ubody.Username)
	log.Printf("Login pwd: %s", pwd)
	log.Printf("Login body pwd: %s", ubody.Pwd)
	if err != nil || len(pwd) == 0 || pwd != ubody.Pwd {
		SendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}

	id := session.GenerateNewSessionId(ubody.Username)
	si := &defs.SignedIn{true, id}
	if resp, err := json.Marshal(si); err != nil {
		SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		SendNormalResponse(w, string(resp), 200)
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		log.Printf("Unauthorized user\n")
		return
	}

	uname := p.ByName("username")
	u, err := dbops.GetUser(uname)
	if err != nil {
		log.Printf("Error in GetUserInfo: %s", err)
		SendErrorResponse(w, defs.ErrorDBError)
		return
	}

	ui := &defs.UserInfo{u.Id}
	if resp, err := json.Marshal(ui); err != nil {
		SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		SendNormalResponse(w, string(resp), 200)
	}
}

func AddNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		log.Printf("Unauthorized user\n")
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	nvbody := &defs.NewVideo{}
	if err := json.Unmarshal(res, nvbody); err != nil {
		log.Printf("%s", err)
		SendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	vi, err := dbops.AddNewVideo(nvbody.AuthorId, nvbody.Name)
	log.Printf("Author id: %d, name: %s \n", nvbody.AuthorId, nvbody.Name)
	if err != nil {
		log.Printf("Error in AddNewVideo: %s", err)
		SendErrorResponse(w, defs.ErrorDBError)
		return
	}

	if resp, err := json.Marshal(vi); err != nil {
		SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		SendNormalResponse(w, string(resp), 201)
	}
}

func ListAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	uname := p.ByName("username")
	vs, err := dbops.ListVideoInfo(uname, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("Error in ListAllvideos: %s", err)
		SendErrorResponse(w, defs.ErrorDBError)
		return
	}

	vsi := &defs.VideosInfo{Videos: vs}
	if resp, err := json.Marshal(vsi); err != nil {
		SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		SendNormalResponse(w, string(resp), 200)
	}

}

func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	vid := p.ByName("vid-id")
	err := dbops.DeleteVideoInfo(vid)
	if err != nil {
		log.Printf("Error in DeleteVideo: %s", err)
		SendErrorResponse(w, defs.ErrorDBError)
		return
	}

	go utils.SendDeleteVideoRequest(vid)
	SendNormalResponse(w, "", 204)
}

func PostComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)

	cbody := &defs.NewComment{}
	if err := json.Unmarshal(reqBody, cbody); err != nil {
		log.Printf("%s", err)
		SendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	vid := p.ByName("vid-id")
	if err := dbops.AddNewComments(vid, cbody.AuthorId, cbody.Content); err != nil {
		log.Printf("Error in Postcomment: %s", err)
		SendErrorResponse(w, defs.ErrorDBError)
	} else {
		SendNormalResponse(w, "ok", 201)
	}
}

func ShowComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	vid := p.ByName("vid-id")
	cm, err := dbops.ListComments(vid, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("Error in ShowComments: %s", err)
		SendErrorResponse(w, defs.ErrorDBError)
		return
	}

	cms := &defs.Comments{cm}
	if resp, err := json.Marshal(cms); err != nil {
		SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		SendNormalResponse(w, string(resp), 200)
	}

}
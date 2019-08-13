package dbops

import (
	"awesomeProject/api/defs"
	"awesomeProject/api/utils"
	"database/sql"
	"log"
	"time"
)

func AddUserCredential(username string, pwd string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO public.users (username, pwd) VALUES ($1, $2)")
	if err != nil {
		log.Printf("error @ AddUserCredential Prepare: %v", err)
		return err
	}

	_, err = stmtIns.Exec(username, pwd)
	if err != nil {
		log.Printf("error @ AddUserCredential Exec: %v", err)
		return err
	}

	defer stmtIns.Close()
	return nil
}

func GetUserCredential(username string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT pwd FROM users WHERE username = $1")
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	var pwd string
	err = stmtOut.QueryRow(username).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	defer stmtOut.Close()

	return pwd, nil
}

func DeleteUser(username string, pwd string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROm users WHERE username = $1 AND pwd = $2")
	if err != nil {
		log.Printf("DeleteUser error: %s", err)
		return err
	}

	_, err = stmtDel.Exec(username, pwd)
	if err != nil {
		return err
	}

	defer stmtDel.Close()
	return nil
}

func GetUser(username string) (*defs.User, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, pwd FROM users WHERE username = $1")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	var id int
	var pwd string

	err = stmtOut.QueryRow(username).Scan(&id, &pwd)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	res := &defs.User{Id: id, Username: username, Pwd: pwd}

	defer stmtOut.Close()

	return res, nil
}

func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {
	// create uuid
	vid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}

	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05")
	stmtIns, err := dbConn.Prepare(`INSERT INTO video_info 
		(id, author_id, name, display_ctime) VALUES($1, $2, $3, $4)`)
	if err != nil {
		return nil, err
	}

	_, err = stmtIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return nil, err
	}

	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: ctime}

	defer stmtIns.Close()
	return res, nil
}

func GetVideoInfo(vid string) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT author_id, name, display_ctime FROM video_info WHERE id=$1")

	var aid int
	var dct string
	var name string

	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	defer stmtOut.Close()

	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: dct}

	return res, nil
}

func ListVideoInfo(uname string, from, to int) ([]*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELECT video_info.id, video_info.author_id, video_info.name, video_info.display_ctime FROM video_info 
		INNER JOIN users ON video_info.author_id = users.id
		WHERE users.username = $1 
		ORDER BY video_info.create_time DESC`)

	var res []*defs.VideoInfo

	if err != nil {
		return res, err
	}

	rows, err := stmtOut.Query(uname)
	if err != nil {
		log.Printf("%s", err)
		return res, err
	}

	for rows.Next() {
		var id, name, ctime string
		var aid int
		if err := rows.Scan(&id, &aid, &name, &ctime); err != nil {
			return res, err
		}

		vi := &defs.VideoInfo{Id: id, AuthorId: aid, Name: name, DisplayCtime: ctime}
		res = append(res, vi)
	}

	defer stmtOut.Close()

	return res, nil
}

func DeleteVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_info WHERE id = $1")
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}

	defer stmtDel.Close()
	return nil
}

func AddNewComments(vid string, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		return err
	}

	stmtIns, err := dbConn.Prepare("INSERT INTO comments (id, video_id, author_id, content) values ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(id, vid, aid, content)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

func ListComments(vid string, from, to int) ([]*defs.Comment, error) {
	stmtOut, err := dbConn.Prepare(` SELECT comments.id, users.username, comments.content FROM comments
		INNER JOIN users ON comments.author_id = users.id
		WHERE comments.video_id = $1
		ORDER BY comments.time DESC`)

	var res []*defs.Comment

	rows, err := stmtOut.Query(vid)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return res, err
		}

		c := &defs.Comment{Id: id, VideoId: vid, Author: name, Content: content}
		res = append(res, c)
	}
	defer stmtOut.Close()

	return res, nil
}

// func ListComments(vid string, from, to int) ([]*defs.Comment, error) {
// 	stmtOut, err := dbConn.Prepare(`SELECT comments.id, users.username, comments.content FROM comments
// 		INNER JOIN users ON comments.author_id = users.id
// 		WHERE comments.video_id = ? AND comments.time > FROM_UNIXTIME(?) AND comments.time <= FROM_UNIXTIME(?)
// 		ORDER BY comments.time DESC`)

// 	var res []*defs.Comment

// 	rows, err := stmtOut.Query(vid, from, to)
// 	if err != nil {
// 		log.Printf("%s", err)
// 		return res, err
// 	}

// 	for rows.Next() {
// 		var id, name, content string
// 		if err := rows.Scan(&id, &name, &content); err != nil {
// 			return res, err
// 		}

// 		c := &defs.Comment{Id: id, VideoId: vid, Author: name, Content: content}
// 		res = append(res, c)
// 	}

// 	defer stmtOut.Close()

// 	return res, nil
// }

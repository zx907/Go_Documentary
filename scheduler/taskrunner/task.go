package taskrunner

import (
	"awesomeProject/scheduler/dbops"
	"errors"
	"log"
	"os"
	"sync"
)

func DeleteVideo(vid string) error {
	err := os.Remove(VIDEO_PATH + "vid")
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Deleting video error: %v", err)
	}
	return nil
}

func VideoClearDispatcher(dc DataChan) error {
	res, err := dbops.ReadVideoDeletionRecord(3)
	if err != nil {
		log.Printf("Video clear dispatcher error: %v", err)
		return err
	}

	if len(res) == 0 {
		return errors.New("all task finished")
	}

	for _, id := range res {
		dc <- id
	}

	return nil
}

func VideoClearExecutor(dc DataChan) error {
	errMap := &sync.Map{}
	var err error
forloop:
	for {
		select {
		case vid := <-dc:
			go func(id interface{}) {
				if err := DeleteVideo(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
				if err := dbops.DelVideoDeletionRecord(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
			}(vid)
		default:
			break forloop
		}
	}

	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil {
			return false
		}
		return true
	})

	return err
}

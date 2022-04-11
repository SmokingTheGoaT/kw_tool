package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hibiken/asynq"
	"time"
)

const (
	TypeRecursiveCrawlRequest = "crawl:recursive"
)

type suggestInterface interface {
	Run() (err error)
}

type RecursiveCrawlPayload struct {
	Request interface{}
}

func NewRecursiveCrawlTask(req interface{}, d time.Duration) (*asynq.Task, error) {
	var task *asynq.Task
	var payload []byte
	var err error
	if _, ok := req.(suggestInterface); ok {
		if payload, err = json.Marshal(RecursiveCrawlPayload{Request: req}); err != nil {
			task = nil
		} else {
			task = asynq.NewTask(TypeRecursiveCrawlRequest, payload, asynq.MaxRetry(1), asynq.ProcessIn(d))
		}
	} else {
		task = nil
		err = errors.New("request has to be of type *requests.Suggest")
	}
	return task, err
}

func HandleRecursiveCrawlTask(ctx context.Context, t *asynq.Task) error {
	var p RecursiveCrawlPayload
	var err error
	if err = json.Unmarshal(t.Payload(), &p); err == nil {
		err = p.Request.(suggestInterface).Run()
	}
	return err
}

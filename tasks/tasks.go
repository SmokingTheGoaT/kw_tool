package tasks

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"kw_tool/types/requests"
	"time"
)

const (
	TypeRecursiveCrawlRequest = "crawl:recursive"
)

type RecursiveCrawlPayload struct {
	Request *requests.Suggest
}

func NewRecursiveCrawlTask(req *requests.Suggest, d time.Duration) (*asynq.Task, error) {
	var task *asynq.Task
	var payload []byte
	var err error
	if payload, err = json.Marshal(RecursiveCrawlPayload{Request: req}); err != nil {
		task = nil
	} else {
		task = asynq.NewTask(TypeRecursiveCrawlRequest, payload, asynq.MaxRetry(1), asynq.ProcessIn(d))
	}
	return task, err
}

func HandleRecursiveCrawlTask(ctx context.Context, t *asynq.Task) error {
	var p RecursiveCrawlPayload
	var err error
	if err = json.Unmarshal(t.Payload(), &p); err == nil {
		err = p.Request.Run()
	}
	return err
}

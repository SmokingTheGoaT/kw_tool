package asynq

import (
	"fmt"
	"github.com/hibiken/asynq"
	"kw_tool/tasks"
	"log"
	"time"
)

type Client struct {
	Cli *asynq.Client
}

func (c *Client) Init(opts asynq.RedisClientOpt) {
	c.Cli = asynq.NewClient(opts)
}

func (c *Client) Enqueue(r interface{}, t time.Duration) (err error) {
	var task *asynq.Task
	if task, err = tasks.NewRecursiveCrawlTask(r, t); err == nil {
		var info *asynq.TaskInfo
		if info, err = c.Cli.Enqueue(task); err == nil {
			log.Println(fmt.Sprintf("enqueued task: id=%s queue=%s", info.ID, info.Queue))
		}
	}
	return
}

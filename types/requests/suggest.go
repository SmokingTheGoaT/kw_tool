package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"kw_tool/tasks"
	"kw_tool/types/crawler"
	"kw_tool/types/keyword"
	"kw_tool/util/common"
	"kw_tool/util/constants"
	"kw_tool/util/enums"
	"kw_tool/util/queue"
	"log"
	"mime"
	"net/http"
	"os"
	"time"
)

type (
	Suggest struct {
		sessionID string
		body      *body
		cfg       *config
	}

	body struct {
		Keywords           string `json:"keywords,omitempty"`
		RecursiveIteration bool   `json:"recursive_iteration,omitempty"`
		PaginationCursor   int    `json:"pagination_cursor,omitempty"`
	}

	config struct {
		client    *asynq.Client
		crawler   *crawler.Crawler
		cache     *cache.Cache
		queue     *queue.Queue
		kwMap     *keyword.Map
		qLength   int
		iteration int
	}

	Response struct {
	}
)

//Init will initialize the request. It takes a cache pointer, time duration and http client as parameters
func (r *Suggest) Init(c *cache.Cache, cli *asynq.Client) {
	r.cfg.iteration = 0 //can be part of body so if passed down in request initial body it will map it.
	r.cfg = &config{
		cache:   c,
		queue:   queue.New(),
		crawler: crawler.New(),
		client:  cli,
	}
}

//ValidateRequest should be initially used on every new request. This method will validate request headers and body
//params
func (r *Suggest) ValidateRequest(req *http.Request) (err error) {
	if err = r.validateHeaders(req); err == nil {
		if err = r.validateBody(req); err == nil {
			r.crawlerConfigurations()
		}
	}
	return
}

//Run will Pop the front item and perform a crawl
func (r *Suggest) Run() (err error) {
	if !(r.cfg.queue.Len() > 0) {
		err = errors.New("")
	} else {
		s := r.cfg.queue.PopFront()
		if str, ok := s.(string); ok {
			str = common.NormalizeQuery(str)
			if err = r.cfg.crawler.Crawl(str); err == nil {
				if err = r.read(str); err == nil {
					if err = r.delete(str); err == nil {
						if r.cfg.iteration < r.cfg.qLength {
							err = r.enqueue()
							r.cfg.iteration++
						} else {
							err = r.save()
						}
					}
				}
			}
		}
	}
	return
}

func (r *Suggest) enqueue() (err error) {
	if r.cfg.iteration == 0 {
		err = r.task(time.Duration(0))
	} else {
		err = r.task(common.GetNextTaskExecutionTime(9000, 18000))
	}
	return
}

func (r *Suggest) task(t time.Duration) (err error) {
	var task *asynq.Task
	if task, err = tasks.NewRecursiveCrawlTask(r, t); err == nil {
		var info *asynq.TaskInfo
		if info, err = r.cfg.client.Enqueue(task); err == nil {
			log.Println(fmt.Sprintf("enqueued task: id=%s queue=%s", info.ID, info.Queue))
		}
	}
	return
}

//save will save cache keyword map to cache instance
func (r *Suggest) save() (err error) {
	r.cfg.cache.Set(r.sessionID, r.cfg.kwMap, cache.NoExpiration)
	if _, found := r.cfg.cache.Get(r.sessionID); !found {
		err = errors.New("error while trying to save keyword map into cache")
	}
	return
}

//validateHeaders will parse headers looking for X-Session-ID and if not found then will create one and return false
//without an error.
func (r *Suggest) validateHeaders(req *http.Request) (err error) {
	//rewrite function to load header to header struct so that it can be used in validateBody function for logic
	var mt string
	if sessID := req.Header.Get(constants.HeaderSessionID); sessID == constants.EmptyString {
		if platform := req.Header.Get(constants.HeaderPlatform); platform == constants.EmptyString {
			err = errors.New("if no X-Session-ID is provided you need to provide a X-Platform header")
		} else {
			r.sessionID = uuid.New().String()
			if mt, _, err = mime.ParseMediaType(platform); err == nil {
				r.cfg.kwMap = &keyword.Map{}
				r.cfg.kwMap.Init(enums.ToPlatform(mt))
			}
		}
	} else {
		if mt, _, err = mime.ParseMediaType(sessID); err == nil {
			r.sessionID = mt
			if km, ok := r.cfg.cache.Get(sessID); ok {
				r.cfg.kwMap = km.(*keyword.Map)
			}
		}
	}
	return
}

//validateBody will parse request body and validate any required ones
func (r *Suggest) validateBody(req *http.Request) (err error) {
	var b []byte
	if b, err = ioutil.ReadAll(req.Body); err == nil {
		var bdy body
		if err = json.Unmarshal(b, &bdy); err == nil {
			if len(r.body.Keywords) == 0 && r.body.PaginationCursor == 0 {
				err = errors.New("there should be at least one keyword sent within request")
			} else if len(r.body.Keywords) > 100 {
				err = errors.New("maximum queries should be 100 or less")
			} else if r.body.PaginationCursor > 0 && len(r.body.Keywords) == 0 {
				// Get the next 100 from keyword map and add it to queue
			} else {
				r.body = &bdy
				r.cfg.qLength = len(r.body.Keywords)
				for _, v := range r.body.Keywords {
					r.cfg.queue.PushBack(v)
				}
			}
		}
	}
	return
}

func (r *Suggest) crawlerConfigurations() {
	p := r.cfg.kwMap.Platform()
	r.cfg.crawler.SetURI(p.URI())
	r.cfg.crawler.SetOpts(p.Options())
}

//read it will check compare against initial length and add to queue accordingly, else it will save results
//to keyword map
func (r *Suggest) read(s string) (err error) {
	var b []byte
	path := fmt.Sprintf("/downloads/%s/%s.json", r.cfg.kwMap.Platform(), s)
	if path, err = common.CreateFilePath(path); err == nil {
		var f *os.File
		if f, err = os.Open(path); err == nil {
			defer func(f *os.File) {
				err = f.Close()
				if err != nil {
					log.Println(err)
				}
			}(f)
			if b, err = ioutil.ReadAll(f); err == nil {
				var re []interface{}
				if err = json.Unmarshal(b, &re); err == nil {
					for i, v := range re {
						if i == 1 {
							if arr, ok := v.([]string); ok {
								for _, x := range arr {
									if r.body.RecursiveIteration {
										if r.cfg.qLength < constants.DefaultSearchCountPerRequest {
											r.cfg.queue.PushBack(x)
											r.cfg.qLength++
										}
									} else {
										r.cfg.kwMap.Insert(x)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

func (r *Suggest) delete(s string) (err error) {
	path := fmt.Sprintf("/downloads/%s/%s.json", r.cfg.kwMap.Platform(), s)
	err = os.Remove(path)
	return
}

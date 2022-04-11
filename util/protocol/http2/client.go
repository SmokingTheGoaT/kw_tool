package http2

import (
	"io"
	"kw_tool/util/common"
	"log"
	"net/http"
	"os"
	"time"
)

type (
	client struct {
		cl http.Client
	}

	ClientInterface interface {
		//DownloadFile is used to copy txt content from a platforms autocomplete api response
		DownloadFile(url string, dest string) (err error)

		//GET(url string, body interface{}, headers map[string][]string) (*http.Response, error)
		//POST(url string, body interface{}, headers map[string][]string) (*http.Response, error)
	}
)

func NewClient() ClientInterface {
	return &client{
		cl: http.Client{
			Timeout: time.Duration(10 * time.Second),
		},
	}
}

//func (c *client) GET(url string, body interface{}, headers map[string][]string) (*http.Response, error) {
//	var err error
//	var rb []byte
//	if rb, err = json.Marshal(body); err == nil {
//		var request *http.Request
//		if request, err = http.NewRequest(http.MethodGet, url, bytes.NewBuffer(rb)); err == nil {
//			request.Header = headers
//			var resp *http.Response
//			if resp, err = c.cl.Do(request); err == nil {
//				return resp, nil
//			}
//		}
//	}
//	return nil, err
//}
//
//func (c *client) POST(url string, body interface{}, headers map[string][]string) (*http.Response, error) {
//	var err error
//	var rb []byte
//	if rb, err = json.Marshal(body); err == nil {
//		var request *http.Request
//		if request, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(rb)); err == nil {
//			request.Header = headers
//			var resp *http.Response
//			if resp, err = c.cl.Do(request); err == nil {
//				return resp, nil
//			}
//		}
//	}
//	return nil, err
//}

//DownloadFile is used to copy txt content from a platforms autocomplete api response
func (c *client) DownloadFile(url string, dest string) (err error) {
	var resp *http.Response
	if resp, err = http.Get(url); err == nil {
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				log.Println(err)
			}
		}(resp.Body)
		if dest, err = common.CreateFilePath(dest); err == nil {
			var out *os.File
			if out, err = os.Create(dest); err == nil {
				defer func(out *os.File) {
					err = out.Close()
					if err != nil {
						log.Println(err)
					}
				}(out)
				_, err = io.Copy(out, resp.Body)
			}
		}
	}
	return
}

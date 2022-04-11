package crawler

import (
	"errors"
	"fmt"
	"kw_tool/util/common"
	"kw_tool/util/protocol/http2"
)

type Options struct {
	Map interface{}
}

func (opt *Options) String() (str string) {
	str = common.ConcatenateMapIntoString(opt.Map.(map[string]string))
	return
}

type Crawler struct {
	opts *Options
	uri  string
	cli  http2.ClientInterface
}

func New() *Crawler {
	return &Crawler{
		opts: &Options{},
		cli:  http2.NewClient(),
	}
}

func (c *Crawler) SetOpts(i map[string]string) {
	c.opts.Map = i
}

func (c *Crawler) SetURI(uri string) {
	c.uri = uri
}

func (c *Crawler) Crawl(q string) (err error) {
	if err = c.prepareRequest(q); err == nil {
		err = c.cli.DownloadFile(c.uri, q)
	}
	return
}

func (c *Crawler) prepareRequest(q string) (err error) {
	if c.uri == "" {
		err = errors.New("uri is empty, please user SetURI to set a uri to crawl")
	} else {
		c.uri = c.uri + c.opts.String() + fmt.Sprintf("&q=%s", q)
	}
	return
}

package network

import (
	pub "common/public_yfm"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"time"
)

type HttpClientHelper struct {
	m2hmsg chan func() func()
	h2mmsg chan func()
	run    bool
}

func (h *HttpClientHelper) process() {
	for i := 0; i < 100; i++ {
		select {
		case f := <-h.m2hmsg:
			cb := f()
			if cb != nil {
				select {
				case h.h2mmsg <- cb:
				default:
					pub.PrintLog("HttpClientHelper h2m msgqueue full.")
					return
				}
			}
		default:
			return
		}
	}
}

func (h *HttpClientHelper) Init() {
	h.run = true
	h.m2hmsg = make(chan func() func(), 1024)
	h.h2mmsg = make(chan func(), 1024)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				pub.PrintLog("HttpClientHelper Work Tread Exit With Error : ", err, " stack : ", debug.Stack())
			}
		}()
		var mili_10 = 10 * time.Millisecond

		for h.run {

			h.process()

			time.Sleep(mili_10)
		}
	}()
}

func (h *HttpClientHelper) Update() {
	for i := 0; i < 100; i++ {
		select {
		case cb := <-h.h2mmsg:
			cb()
		default:
			return
		}
	}
}

func (h *HttpClientHelper) HttpPost(url, contentType string, body io.Reader, callback func([]byte)) bool {
	task := func() func() {
		if rsp, err := http.Post(url, contentType, body); err != nil {
			if res, ioerr := ioutil.ReadAll(body); ioerr != nil {
				pub.PrintLog("HttpClientHelper HttpPost Fail with err : ", err.Error(), " url : ", url, "body : ioread err : ", ioerr.Error())
			} else {
				pub.PrintLog("HttpClientHelper HttpPost Fail with err : ", err.Error(), " url : ", url, "body : ", string(res))
			}
			return nil
		} else {
			if res, ioerr := ioutil.ReadAll(rsp.Body); ioerr != nil {
				pub.PrintLog("HttpClientHelper HttpPost readfrom rsp err : ", ioerr.Error(), " url : ", url)
				return nil
			} else {
				if callback != nil {
					return func() {
						callback(res)
					}
				} else {
					return nil
				}
			}
		}
	}

	select {
	case h.m2hmsg <- task:
		return true
	default:
		pub.PrintLog("HttpClientHelper m2h msgqueue full.")
		return false
	}
}

func (h *HttpClientHelper) HttpGet(url string, callback func([]byte)) bool {
	task := func() func() {
		if rsp, err := http.Get(url); err != nil {
			pub.PrintLog("HttpClientHelper HttpGet Fail with err : ", err.Error(), " url : ", url)
			return nil
		} else {
			if res, ioerr := ioutil.ReadAll(rsp.Body); ioerr != nil {
				pub.PrintLog("HttpClientHelper readfrom rsp err : ", ioerr.Error(), " url : ", url)
				return nil
			} else {
				if callback != nil {
					return func() {
						callback(res)
					}
				} else {
					return nil
				}
			}

		}
	}

	select {
	case h.m2hmsg <- task:
		return true
	default:
		pub.PrintLog("HttpClientHelper m2h msgqueue full.")
		return false
	}
}

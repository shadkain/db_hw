package delivery

import (
	"github.com/valyala/fasthttp"
)

func (this *Handler) getStatus(c *fasthttp.RequestCtx) {
	if status, err := this.uc.GetStatus(); err != nil {
		Error(c, err)
	} else {
		Ok(c, status)
	}
}

func (this *Handler) clear(c *fasthttp.RequestCtx) {
	if err := this.uc.Clear(); err != nil {
		Error(c, err)
	} else {
		Ok(c, nil)
	}
}

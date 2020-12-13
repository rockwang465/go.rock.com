package v1

import "rock-platform/rock/server/log"

type Controller struct {
	*log.Logger
}

var ctl *Controller

func GetController() *Controller {
	if ctl == nil {
		ctl = &Controller{
			log.GetLogger(),
		}
	}
	return ctl
}

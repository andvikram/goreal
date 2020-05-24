package controller

import (
	"github.com/andvikram/goreal/logger"
)

var (
	log = logger.GoRealLog{}
	// RouteMap contains API paths
	RouteMap = make(map[string]string)
)

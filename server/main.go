package main

import (
	"SpectreLink/admin"
	"SpectreLink/log"
	glog "github.com/charmbracelet/log"
)

func main() {
	log.SetLevel(glog.DebugLevel)
	spec := admin.NewSpectreLink()
	spec.ListenAndServe()
	
}

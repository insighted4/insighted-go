package main

import (
	"github.com/insighted4/insighted-go/examples/github/api"
	"github.com/insighted4/insighted-go/kit"
)

func main() {
	cfg := kit.DefaultConfig()
	cfg.LoggerFormat = "text"
	cfg.EnablePProf = true
	svc := api.New(cfg)
	kit.Run(svc)
}

package main

import (
	"github.com/insighted4/insighted-go/kit"
	"github.com/insighted4/insighted-go/examples/github/api"
)

func main() {
	cfg := kit.DefaultConfig()
	cfg.LoggerFormat = "text"
	cfg.EnablePProf = true
	svc := api.New()
	kit.Run(cfg, svc, )
}

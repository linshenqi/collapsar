package main

import (
	"flag"
	"github.com/linshenqi/collapsar/src/services/oss"
	"github.com/linshenqi/sptty"
)

func main() {
	cfg := flag.String("config", "./config.yml", "--config")
	flag.Parse()

	app := sptty.GetApp()
	app.ConfFromFile(*cfg)

	services := sptty.Services{
		&oss.Service{},
	}

	configs := sptty.Configs{
		&oss.Config{},
	}

	app.AddServices(services)
	app.AddConfigs(configs)

	app.Sptting()
}

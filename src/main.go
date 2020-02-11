package main

import (
	"flag"
	"github.com/linshenqi/collapsar/src/services/oss"
	"github.com/linshenqi/collapsar/src/services/qiniu"
	"github.com/linshenqi/sptty"
)

func main() {
	cfg := flag.String("config", "./config.yml", "--config")
	flag.Parse()

	app := sptty.GetApp()
	app.ConfFromFile(*cfg)

	ossService := oss.Service{}
	ossService.SetupProviders(map[string]oss.IOss{
		oss.Qiniu: &qiniu.Oss{},
	})

	services := sptty.Services{
		&ossService,
	}

	configs := sptty.Configs{
		&oss.Config{},
	}

	app.AddServices(services)
	app.AddConfigs(configs)

	app.Sptting()
}

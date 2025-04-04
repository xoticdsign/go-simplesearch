package app

import (
	"github.com/xoticdsign/simplesearch/internal/app/httpsss"
	"github.com/xoticdsign/simplesearch/internal/lib/logger"
	"github.com/xoticdsign/simplesearch/internal/utils"
)

type App struct {
	SimpleSearch *httpsss.App
}

func New(env string) (App, error) {
	cfg, err := utils.MustLoadConfig(env)
	if err != nil {
		return App{}, err
	}

	log := logger.New(env)

	ss := httpsss.New(log, cfg)

	return App{
		SimpleSearch: ss,
	}, nil
}

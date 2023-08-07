package probes

import (
	"auth/internal/config"
	"auth/pkg/infra/logger"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Probes struct {
	isReady   bool
	isStarted bool

	readyOnce   sync.Once
	startedOnce sync.Once

	logger logger.Logger
}

func New(logger logger.Logger) (*Probes, error) {
	return &Probes{
		logger: logger,
	}, nil
}

func (p *Probes) Start() error {
	cfg, _ := config.GetConfig()

	r := gin.Default()
	r.GET("/healthz", func(ctx *gin.Context) {
		ctx.Writer.WriteHeader(http.StatusOK)
	})

	r.GET("/ready", func(ctx *gin.Context) {
		if p.isReady {
			ctx.Writer.WriteHeader(http.StatusOK)
		} else {
			ctx.Writer.WriteHeader(http.StatusLocked)
		}
	})

	r.GET("/startup", func(ctx *gin.Context) {
		if p.isStarted {
			ctx.Writer.WriteHeader(http.StatusOK)
		} else {
			ctx.Writer.WriteHeader(http.StatusLocked)
		}
	})

	go func() {
		err := r.Run(":" + cfg.ProbesPort)
		if err != nil {
			p.logger.Sugar().Errorf("start probes failed: %s", err.Error())
		}
	}()

	return nil
}

func (p *Probes) SetReady() {
	p.readyOnce.Do(func() {
		p.isReady = true
	})
}

func (p *Probes) SetStarted() {
	p.startedOnce.Do(func() {
		p.isStarted = true
	})
}

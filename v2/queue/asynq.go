package queue

import (
	"context"
	kc "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hibiken/asynq"
)

type RedisConfig struct {
	Network  string `json:"network"`
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func NewAsynqProducer(kc kc.Config, logger log.Logger) (Producer, func()) {
	var opt RedisConfig
	err := kc.Value("data.redis").Scan(&opt)
	if err != nil {
		panic(err)
	}
	client := asynq.NewClient(asynq.RedisClientOpt{
		Network:  opt.Network,
		Addr:     opt.Addr,
		Username: opt.Username,
		Password: opt.Password,
		DB:       opt.DB,
	})

	return &AsynqProducer{
			client,
			log.NewHelper(log.With(logger, "module", "asynq-producer")),
		}, func() {
			client.Close()
		}
}

type AsynqProducer struct {
	client *asynq.Client
	logger *log.Helper
}

func (p AsynqProducer) Send(message Message) error {
	p.logger.Debugf("send message topic: %s, payload: %s", message.Topic, string(message.Payload))
	_, err := p.client.Enqueue(asynq.NewTask(message.Topic, message.Payload))
	if err != nil {
		p.logger.Debugf("failed to enqueue: %v", err)
	}
	return err
}

type AsynqConsumers struct {
	server   *asynq.Server
	logger   *log.Helper
	handlers map[string]func(Message) error
}

func NewAsynqConsumers(kc kc.Config, logger log.Logger) Consumers {
	var opt RedisConfig
	err := kc.Value("data.redis").Scan(&opt)
	if err != nil {
		panic(err)
	}
	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Network:  opt.Network,
			Addr:     opt.Addr,
			Username: opt.Username,
			Password: opt.Password,
			DB:       opt.DB,
		},
		asynq.Config{
			Concurrency: 16,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)
	return &AsynqConsumers{
		server:   server,
		logger:   log.NewHelper(log.With(logger, "module", "asynq-consumers")),
		handlers: map[string]func(Message) error{},
	}
}

func (c AsynqConsumers) RegisterHandler(topic string, handler func(Message) error) {
	c.handlers[topic] = handler
}

func (c AsynqConsumers) Start() error {
	mux := asynq.NewServeMux()
	for topic, handler := range c.handlers {
		mux.HandleFunc(topic, func(ctx context.Context, t *asynq.Task) error {
			return handler(Message{
				Topic:   t.Type(),
				Payload: t.Payload(),
			})
		})
	}
	err := c.server.Start(mux)
	if err != nil {
		c.logger.Errorf("failed to start asynq server: %v", err)
	} else {
		c.logger.Infof("asynq server started")
	}
	return err
}

func (c AsynqConsumers) Stop() error {
	c.server.Stop()
	c.server.Shutdown()
	c.logger.Infof("asynq server stopped")
	return nil
}

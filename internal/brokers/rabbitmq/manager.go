package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tolikproh/banners-rotation/internal/config"
	"github.com/tolikproh/banners-rotation/internal/logger"
	"github.com/tolikproh/banners-rotation/internal/model"
)

type Rabbit struct {
	cfg         *config.Config
	log         *logger.Logger
	consumerTag string
	channel     *amqp.Channel
}

func NewRabbit(ctx context.Context, cfg *config.Config, log *logger.Logger) (*Rabbit, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.Address)
	if err != nil {
		return nil, fmt.Errorf("error connect to rabbit (%s): %w", cfg.RabbitMQ.Address, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error open rabbit channel (%s): %w", cfg.RabbitMQ.Address, err)
	}

	if len(cfg.RabbitMQ.Exchange) > 0 {
		if err = ch.ExchangeDeclare(
			cfg.RabbitMQ.Exchange,
			amqp.ExchangeDirect,
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return nil, fmt.Errorf("error declare exchange (%s): %w", cfg.RabbitMQ.Exchange, err)
		}
	}

	q, err := ch.QueueDeclare(
		cfg.RabbitMQ.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error declare queue (%s): %w", cfg.RabbitMQ.Queue, err)
	}

	if err = ch.QueueBind(
		q.Name,
		q.Name,
		cfg.RabbitMQ.Exchange,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("error bind queue: %w", err)
	}

	go func() {
		<-ctx.Done()
		ch.Close()
		conn.Close()
	}()

	return &Rabbit{
		cfg:         cfg,
		log:         log,
		consumerTag: "banners-consumer",
		channel:     ch,
	}, nil
}

func (q *Rabbit) Add(n model.Event) error {
	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("error marshall StatEvent: %w", err)
	}

	if err = q.channel.Publish(
		q.cfg.RabbitMQ.Exchange,
		q.cfg.RabbitMQ.Queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		}); err != nil {
		return fmt.Errorf("error publish StatEvent: %w", err)
	}

	return nil
}

package nats

import (
	"context"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

const otelName = "b0b-common/nats"

type Client struct {
	Conn *nats.Conn
}
type JetClient struct {
	*Client
	Conn nats.JetStreamContext
}

func NewClient() (*Client, error) {
	nc, err := nats.Connect(config.Cfg.NatsCfg.Host, nats.UserInfo(config.Cfg.NatsCfg.Username, config.Cfg.NatsCfg.Password))
	if err != nil {
		return nil, errors.Wrap(err, "connect fail")
	}
	return &Client{Conn: nc}, nil
}

func NewJetClient() (*JetClient, error) {
	nc, err := nats.Connect(config.Cfg.NatsCfg.Host, nats.UserInfo(config.Cfg.NatsCfg.Username, config.Cfg.NatsCfg.Password))
	if err != nil {
		return nil, errors.Wrap(err, "connect fail")
	}
	stream, err := nc.JetStream()
	if err != nil {
		return nil, errors.Wrap(err, "JetStream fail")
	}
	return &JetClient{Conn: stream, Client: &Client{Conn: nc}}, nil
}

func (c *Client) Publish(ctx context.Context, subj string, data []byte) error {
	defer newOTELSpan(ctx, "MQ.Publish").End()
	return c.Conn.Publish(subj, data)
}
func (c *JetClient) Publish(ctx context.Context, subj string, data []byte) (*nats.PubAck, error) {
	defer newOTELSpan(ctx, "MQ.JetPublish").End()
	return c.Conn.Publish(subj, data)
}

func newOTELSpan(ctx context.Context, name string) trace.Span {
	_, span := otel.Tracer(otelName).Start(ctx, name)
	span.SetAttributes(attribute.KeyValue{Key: semconv.MessagingSystemKey, Value: attribute.StringValue("nats")})
	return span
}

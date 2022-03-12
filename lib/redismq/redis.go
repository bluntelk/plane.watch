package redismq

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"net"
	"net/url"
	"time"
)

type Server struct {
	url      string
	incoming *redis.Client
	pubSub   *redis.PubSub
	outgoing *redis.Client
}

func NewServer(serverUrl string) (*Server, error) {
	r := &Server{}
	r.SetUrl(serverUrl)
	if err := r.Connect(); nil != err {
		return nil, err
	}
	return r, nil
}

func (r *Server) SetUrl(serverUrl string) {
	r.url = serverUrl
}

func (r *Server) Connect() error {
	var err error
	serverUrl, err := url.Parse(r.url)
	if nil != err {
		return err
	}
	log.Debug().Str("url", r.url).Msg("connecting to redis server...")
	pass, _ := serverUrl.User.Password()
	port := serverUrl.Port()
	if "" == port {
		port = "6379"
	}
	opts := &redis.Options{
		Addr:     net.JoinHostPort(serverUrl.Hostname(), port),
		Username: serverUrl.User.Username(),
		Password: pass,
		DB:       0,
	}

	r.incoming = redis.NewClient(opts)
	r.outgoing = redis.NewClient(opts)
	return nil
}

// Publish is our simple message publisher
func (r *Server) Publish(queue string, msg []byte) error {
	ctx := context.Background()
	if out := r.outgoing.Publish(ctx, queue, msg); nil != out {
		return out.Err()
	}
	return nil
}

func (r *Server) Close() {
	if nil != r.pubSub {
		_ = r.pubSub.Close()
	}
	if nil != r.incoming {
		_ = r.incoming.Close()
	}
	if nil != r.outgoing {
		_ = r.outgoing.Close()
	}
}

func (r *Server) Subscribe(subject string) (<-chan *redis.Message, error) {
	ctx := context.Background()
	r.pubSub = r.incoming.Subscribe(ctx, subject)

	ch := r.pubSub.Channel()
	return ch, nil
}

func (r *Server) HealthCheckName() string {
	return "Redis PubSub"
}

func (r *Server) HealthCheck() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	in := r.incoming.Ping(ctx)
	out := r.incoming.Ping(ctx)
	return (nil != in && in.Err() != nil) || (nil != out && out.Err() != nil)
}

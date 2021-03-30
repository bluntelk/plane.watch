package sink

import (
	"io"
	"os"
	"sync"
)

type (
	Config struct {
		host, port string
		secure bool
		queue string

		out io.WriteCloser
		waiter sync.WaitGroup
	}
	Option func(*Config)
)

func WithHost(host, port string) Option {
	return func(conf *Config) {
		conf.host = host
		conf.port = port
	}
}

func WithQueue(queue string) Option {
	return func(conf *Config) {
		conf.queue = queue
	}
}

func WithLogOutput(out io.WriteCloser) Option {
	return func(config *Config) {
		config.out = out
	}
}

func WithLogFile(file string) Option {
	return func(config *Config) {
		f, err := os.Create(file)
		if nil != err {
			println("Cannot open file: ", file)
			return
		}
		config.out = f
	}
}

func (c *Config) Finish() {
	c.waiter.Wait()
	_ = c.out.Close()
}
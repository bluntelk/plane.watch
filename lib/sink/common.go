package sink

import (
	"io"
	"os"
)

type (
	Config struct {
		host, port string
		secure bool
		queue string

		out io.Writer
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

func WithLogOutput(out io.Writer) Option {
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

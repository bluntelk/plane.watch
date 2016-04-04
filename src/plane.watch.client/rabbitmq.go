package main

type RabbitMqSender struct {
	host, port string
}

func NewRabbitSender(host, port string) *RabbitMqSender {
	rmq := new(RabbitMqSender)
	rmq.host = host
	rmq.port = port
}
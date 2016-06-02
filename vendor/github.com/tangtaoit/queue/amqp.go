package queue

import (
	"github.com/streadway/amqp"
	"github.com/tangtaoit/util"
)

var connection *amqp.Connection


func SetupAMQP(url string)  {

	if url=="" {
		url ="amqp://guest:guest@localhost:5673/"
	}
	var err error
	connection, err = amqp.Dial("amqp://guest:guest@localhost:5673/")
	util.CheckErr(err)
}

func GetChannel() *amqp.Channel {

	channel, err := connection.Channel()
	util.CheckErr(err)
	return channel
}
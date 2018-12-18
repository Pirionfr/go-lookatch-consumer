package core

import (
	"strings"

	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	log "github.com/sirupsen/logrus"
)

// Consume consumes messages from Kafka
func Consume(brokerAddrs string, topic string, username string, password string, consumerGroup string) <-chan *sarama.ConsumerMessage {

	var config = cluster.NewConfig()
	config.Net.TLS.Enable = true
	config.Net.SASL.Enable = true
	config.Net.SASL.User = username
	config.Net.SASL.Password = password
	config.ClientID = username

	if consumerGroup == "" {
		consumerGroup = username + ".go"
	}

	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	var err error
	consumer, err := cluster.NewConsumer(strings.Split(brokerAddrs, ","), brokerAddrs, []string{topic}, config)
	HandleError(err)

	// Consume errors
	go func() {
		for err := range consumer.Errors() {
			log.WithError(err).Error("Error during consumption")
		}
	}()

	// Consume notifications
	go func() {
		for note := range consumer.Notifications() {
			log.WithField("note", note).Debug("Rebalanced consumer")
		}
	}()

	return consumer.Messages()

}

func HandleError(err error) {
	if err != nil {
		fmt.Println("An error occurred: ", err.Error())
		os.Exit(1)
	}
}

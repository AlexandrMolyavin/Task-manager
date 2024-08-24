package kafka

import (
	"github.com/IBM/sarama"
)

func ConnectConsumer(brokers []string) (sarama.PartitionConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	partitionConsumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}
	consumer, err := partitionConsumer.ConsumePartition("tasks", 0, sarama.OffsetNewest)
	return consumer, err
}

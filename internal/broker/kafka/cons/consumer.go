package cons

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type ConsumerKafka struct {
	Cons sarama.PartitionConsumer
}

func NewConsumer(cons sarama.PartitionConsumer) ConsumerKafka {
	return ConsumerKafka{Cons: cons}
}

func (cn *ConsumerKafka) GetReq(c *gin.Context) []byte {
	var res []byte
	select {
	case err := <-cn.Cons.Errors():
		fmt.Println(err)
	case msg := <-cn.Cons.Messages():
		fmt.Printf("Offset: %v, Key: %s, Value: %s\n", msg.Offset, msg.Key, msg.Value)
		res = msg.Value
	}
	return res
}

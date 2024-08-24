package prod

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"net/http"
	"ts/internal/manager"
)

type ProducerKafka struct {
	prod sarama.SyncProducer
}

func NewProducer(prod sarama.SyncProducer) ProducerKafka {
	return ProducerKafka{prod: prod}
}
func (p *ProducerKafka) PushReqToQueue(c *gin.Context, topic string, body []byte) error {

	// Create new Kafka message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(body),
	}
	//Send message
	_, _, err := p.prod.SendMessage(msg)
	if err != nil {
		return err
	}
	c.JSON(200, gin.H{
		"Info": "",
	})
	return nil
}

func (p *ProducerKafka) PlaceReq(c *gin.Context) {
	// Parse request body
	var (
		taskDecode manager.Task
	)

	err := json.NewDecoder(c.Request.Body).Decode(&taskDecode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to decode request body",
		})
		return
	}
	// convert req body into bytes
	taskInBytes, err := json.Marshal(taskDecode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to convert order into bytes",
		})
		return
	}

	// send bytes to Kafka
	err = p.PushReqToQueue(c, "tasks", taskInBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to push task to queue",
		})
		panic(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"INFO": c.Request.Method + " method successfully pushed to Kafka queue",
	})

}

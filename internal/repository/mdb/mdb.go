package mdb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	manager2 "ts/internal/manager"
)

type MongoDB struct {
	Mongo *mongo.Collection
}

func NewMongoDB(db *mongo.Collection) MongoDB {
	return MongoDB{Mongo: db}
}

func (m *MongoDB) FindFirst(user interface{}, login string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDB) CreateUser(user interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDB) Post(c *gin.Context, bRes []byte, logger *zerolog.Logger) {
	var taskDecode manager2.Task
	err := json.NewDecoder(c.Request.Body).Decode(&taskDecode)
	if err != nil {
		http.Error(c.Writer, "Unsupported method", http.StatusMethodNotAllowed)
	}
	taskDecode.Status = "Proccessing"
	_, err2 := m.Mongo.InsertOne(context.Background(), taskDecode)

	if err2 != nil {
		http.Error(c.Writer, err2.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(c.Writer, "Task was added '%s'", taskDecode.Id)
}

func (m *MongoDB) GetById(c *gin.Context, bRes []byte, logger *zerolog.Logger) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDB) GetAll(c *gin.Context, logger *zerolog.Logger) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDB) ChangeStatus(c *gin.Context, bRes []byte, logger *zerolog.Logger) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDB) DeleteById(c *gin.Context, bRes []byte, logger *zerolog.Logger) {
	//TODO implement me
	panic("implement me")
}

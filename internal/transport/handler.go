package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"ts/internal/broker/kafka/cons"
	"ts/internal/broker/kafka/prod"
)

type Manager interface {
	Post(c *gin.Context, bRes []byte, logger *zerolog.Logger)
	GetById(c *gin.Context, bRes []byte, logger *zerolog.Logger)
	GetAll(c *gin.Context, logger *zerolog.Logger)
	ChangeStatus(c *gin.Context, bRes []byte, logger *zerolog.Logger)
	DeleteById(c *gin.Context, bRes []byte, logger *zerolog.Logger)
}

type Handler struct {
	Repo   Manager
	Logger *zerolog.Logger
	Prod   prod.ProducerKafka
	Cons   cons.ConsumerKafka
}

func NewHandler(repo Manager) *Handler {
	return &Handler{Repo: repo}
}

func (h *Handler) HandleRequest(c *gin.Context) {
	method := c.Request.Method
	h.Prod.PlaceReq(c)
	res := h.Cons.GetReq(c)
	switch method {
	case "GET":
		h.Repo.GetById(c, res, h.Logger)
	case "POST":
		h.Repo.Post(c, res, h.Logger)
	case "DELETE":
		h.Repo.DeleteById(c, res, h.Logger)
	case "PUT":
		h.Repo.ChangeStatus(c, res, h.Logger)
	default:
		http.Error(c.Writer, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

package pdb

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"ts/internal/manager"
)

type Pgx struct {
	Db *gorm.DB
}

func NewPgx(db *gorm.DB) Pgx {
	return Pgx{Db: db}
}

func (p *Pgx) CreateUser(user interface{}) error {
	result := p.Db.Create(user)
	return result.Error
}

func (p *Pgx) FindFirst(user interface{}, scond string, cond interface{}) error {
	var result *gorm.DB
	if scond == "" {
		result = p.Db.First(user, cond)
	} else {
		result = p.Db.First(user, scond, cond)
	}
	return result.Error
}

func (p *Pgx) Post(c *gin.Context, bRes []byte, logger *zerolog.Logger) {
	var taskDecode manager.Task
	err := json.Unmarshal(bRes, &taskDecode)
	//err := json.NewDecoder(c.Request.Body).Decode(&taskDecode)
	if err != nil {
		logger.Error().Msg("Cannot get ID from request")
		fmt.Fprintln(c.Writer, "Cannot get ID from request")
	}
	logger.Info().Str("ID", taskDecode.Id).Msg("POST request captured")
	taskDecode.Status = "Proccessing"
	result := p.Db.Create(taskDecode)
	if result.Error != nil {
		fmt.Println("Error inserting record:", result.Error)
	} else {
		fmt.Println("Record inserted successfully")
	}

	logger.Info().Str("method", c.Request.Method).Str("ID", taskDecode.Id).Msg("Task was added")
	fmt.Fprintf(c.Writer, "Task was added '%s'", taskDecode.Id)
}

func (p *Pgx) GetById(c *gin.Context, bRes []byte, logger *zerolog.Logger) {
	var (
		taskDecode manager.Task
		taskFound  manager.Task
	)
	err := json.Unmarshal(bRes, &taskDecode)
	//err := json.NewDecoder(c.Request.Body).Decode(&taskDecode)
	if err != nil {
		logger.Error().Msg("Cannot get ID from request")
		fmt.Fprintln(c.Writer, "Cannot get ID from request")
	} else {
		logger.Info().Str("ID", taskDecode.Id).Msg("GET by ID request captured")
	}
	if taskDecode.Id == "%ALL" {
		p.GetAll(c, logger)
		return
	}
	errF := p.Db.Where("id = ?", taskDecode.Id).First(&taskFound).Error
	if errF != nil {
		logger.Error().Str("ID", taskDecode.Id).Msg("Cannot get task from DB")
		fmt.Fprintln(c.Writer, "Task not found")
	} else {
		logger.Info().Str("ID", taskFound.Id).Msg("Task found")
		fmt.Fprintf(c.Writer, "Task ID: %s\nStatus: %s", taskFound.Id, taskFound.Status)
	}
}

func (p *Pgx) GetAll(c *gin.Context, logger *zerolog.Logger) {
	var (
		tm = make([]manager.Task, 20)
	)
	result := p.Db.Find(&tm)

	fmt.Fprintf(c.Writer, "%d tasks added \n", len(tm))
	if len(tm) != 0 {
		for _, task := range tm {
			json.NewEncoder(c.Writer).Encode(task)
		}
	}

	if result.Error != nil {
		fmt.Println("Error querying records:", result.Error)
	} else {
		fmt.Printf("Records found: %+v\n", tm)
	}

	logger.Info().Msg("GET All request successful")
}

/*Изменение статуса задачи по ID*/
func (p *Pgx) ChangeStatus(c *gin.Context, bRes []byte, logger *zerolog.Logger) {
	//logger.Info().Msg("PUT request captured")
	var (
		taskDecode manager.Task
	)

	err := json.Unmarshal(bRes, &taskDecode)
	if err != nil {
		logger.Error().Msg("Cannot decode request body")
		fmt.Fprintln(c.Writer, "Cannot decode request body")
		return
	} else {
		logger.Info().Str("ID", taskDecode.Id).Msg("PUT request captured")
	}
	res := p.Db.Model(&taskDecode).Where("id = ?", taskDecode.Id).
		Updates(manager.Task{Id: taskDecode.Id, Status: taskDecode.Status})

	if res.Error != nil {
		fmt.Fprintf(c.Writer, "Err: %v", res.Error)
	}

	if res.RowsAffected == 0 {
		logger.Error().Msg("Cannot update status, task not found")
		fmt.Fprintln(c.Writer, "Cannot update status, task not found")
	} else {
		fmt.Fprintf(c.Writer, "Status was updated:\nTask ID: %s\nStatus: %s", taskDecode.Id, taskDecode.Status)
		logger.Info().Str("ID", taskDecode.Id).Str("Status", taskDecode.Status).Msg("Status was updated")
	}
}

/*Удаление задачи по ID*/
func (p *Pgx) DeleteById(c *gin.Context, bRes []byte, logger *zerolog.Logger) {
	var taskDecode manager.Task

	errD := json.Unmarshal(bRes, &taskDecode)
	if errD != nil {
		logger.Error().Msg("Cannot get ID from request")
		fmt.Fprintf(c.Writer, "Cannot get ID from request")
	} else {
		logger.Info().Str("ID", taskDecode.Id).Msg("DELETE by ID request captured")
	}
	res := p.Db.Where("id = ?", taskDecode.Id).Delete(&taskDecode)

	if res.Error != nil {
		fmt.Fprintf(c.Writer, "Err: %v", res.Error)
	}

	if res.RowsAffected == 0 {
		fmt.Fprintln(c.Writer, "Task not found")
	} else {
		fmt.Fprintf(c.Writer, "Task '%s' was deleted ", taskDecode.Id)
		logger.Info().Str("ID", taskDecode.Id).Msg("Task was deleted")
	}
}

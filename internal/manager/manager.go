package manager

type Task struct {
	Id     string `json:"TaskID" bson:"TaskID"`
	Status string `json:"TaskStatus" bson:"TaskStatus"`
}

type TaskMap map[string]string

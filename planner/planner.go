package planner

import (
	"errors"

	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"log"
)

type Task struct {
	ID            int    `json:"id"`
	Priority      int    `json:"priority"`
	Title_Of_Task string `json:"title_of_task"`
	Done          bool   `json:"done"`
	Order         int    `json:"order"`
}

type Planner struct {
	Tasks  map[int]Task `json:"tasks"`
	NextID int
}

func (p *Planner) GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	task, found := p.Tasks[id]
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (p *Planner) Reset(c *gin.Context) {
    p.Tasks = make(map[int]Task)
    p.NextID = 0
    c.JSON(http.StatusOK, gin.H{"status": "Planner state reset"})
}



func (p *Planner) AddTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !task.Done {
		task.Done = false
	}
	p.NextID++
	task.ID = p.NextID
	task.Priority = len(p.Tasks) + 1
	p.Tasks[task.ID] = task
	c.JSON(http.StatusCreated, task)
}
func (p *Planner) SetPriority(id, priority int) error {

	task, exist := p.Tasks[id]
	if !exist {
		return errors.New("task not found")
	}

	task.Priority = priority
	p.Tasks[id] = task
	return nil

}

// UpdateTaskDone updates the Done status of a task
func (p *Planner) UpdateTaskDone(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	task, found := p.Tasks[id]
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var jsonData struct {
		Done bool `json:"done"`
	}

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.Done = jsonData.Done
	p.Tasks[id] = task

	c.JSON(http.StatusOK, task)
}

func (p *Planner) UpdateTask(c *gin.Context) {
	 log.Println("Entering UpdateTask function")
	id,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		  log.Println("Error parsing ID:", err)
		c.JSON(http.StatusBadRequest,gin.H{"error": "Invalid ID"})
		return
	}

	task,found := p.Tasks[id]
	if !found {
		 log.Println("Task not found for ID:", id)
		c.JSON(http.StatusBadRequest,gin.H{"error": "Task Not Found"})
		return
	}
	var updatedTask Task
		if err := c.ShouldBindJSON(&updatedTask); err != nil {
	   log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
}

	updatedTask.ID = task.ID;

	updatedTask.Priority = task.Priority

	p.Tasks[id] = updatedTask
	c.JSON(http.StatusOK,updatedTask)
}

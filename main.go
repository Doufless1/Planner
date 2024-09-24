package main

import (
	"planner"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the Planner with empty tasks and reset NextID
	p := &planner.Planner{
		Tasks:  make(map[int]planner.Task),
		NextID: 0, // Ensure NextID is reset at the start
	}

	// Log the initialization of the Planner
	p.Reset(nil)

	r:= gin.Default()

// Define your routes
	r.POST("/tasks", p.AddTask)
	r.PUT("/tasks/:id", p.UpdateTask)
	r.GET("/tasks/:id", p.GetTask)
	r.PUT("/tasks/:id/done", p.UpdateTaskDone)
	r.DELETE("/reset", p.Reset)
	r.Run(":8080")
	
}

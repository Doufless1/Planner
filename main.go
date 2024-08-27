package main

import (
	"github.com/gin-gonic/gin"
	"planner"
)

func main() {
	// Initialize the Planner with empty tasks and reset NextID
	p := &planner.Planner{
		Tasks:  make(map[int]planner.Task),
		NextID: 0, // Ensure NextID is reset at the start
	}

	// Log the initialization of the Planner

	router := gin.Default()

	// Define routes and handlers
	router.POST("/tasks",	p.AddTask)

	router.GET("/tasks/:id", p.GetTask)
	//   router.PUT("/tasks/:id/priority", p.SetPriority)

	router.PATCH("/tasks/:id/done", p.UpdateTaskDone)
	router.Run(":8080")
}

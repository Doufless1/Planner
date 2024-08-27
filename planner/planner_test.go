package planner_test

import (
	"planner"
	"testing"

	"net/http"
	"net/http/httptest"

	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	//	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/assert"
)

func TestAddTask( t *testing.T){
	t.Parallel()

	gin.SetMode(gin.TestMode)

	p := &planner.Planner{
		Tasks: make(map[int]planner.Task),
		NextID: 0,
	}

   tests := []struct {
        name           string
        body           interface{}
        expectedStatus int
        expectedBody   string
    }{
        {
            name:           "Valid Task",
            body:           planner.Task{Title_Of_Task: "New Task", Done: false},
            expectedStatus: http.StatusCreated,
            expectedBody:   `{"id":1,"priority":1,"title_of_task":"New Task","done":false}`,
        },
        {
            name:           "Invalid JSON",
            body:           `{"title_of_task":"New Task", "done":"invalid_value"}`, // Invalid 'done' field
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"json: cannot unmarshal string into Go struct field Task.done of type bool"}`,
        },
    }

	for _, tc := range tests {
     t.Run(tc.name,func(t *testing.T) {
			 var bodyBytes []byte
			 var err error
			 switch v:= tc.body.(type){
			 case string:
				 bodyBytes = []byte(v)
			 default:
				 bodyBytes, err = json.Marshal(v)
				 assert.NoError(t,err)
			 }
req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(bodyBytes))
            assert.NoError(t, err)
            req.Header.Set("Content-Type", "application/json")

            rr := httptest.NewRecorder()

            c, _ := gin.CreateTestContext(rr)
            c.Request = req

            p.AddTask(c)

            // Check status code
            assert.Equal(t, tc.expectedStatus, rr.Code)

            // Check body
            assert.JSONEq(t, tc.expectedBody, rr.Body.String())})
	}
	 t.Run("Check Task State", func(t *testing.T) {
        assert.Equal(t, 1, len(p.Tasks))
        addedTask, exists := p.Tasks[1]
        assert.True(t, exists)
        assert.Equal(t, 1, addedTask.ID)
        assert.Equal(t, 1, addedTask.Priority)
        assert.Equal(t, "New Task", addedTask.Title_Of_Task)
        assert.False(t, addedTask.Done)
    })
}

func TestGetTask(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	planner := &planner.Planner{
		Tasks: map[int]planner.Task{
			1: {ID: 1, Priority: 1, Title_Of_Task: "Test", Done: false},
		},
	}

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid ID",
			id:             "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"priority":1,"title_of_task":"Test","done":false}`,
		},
		{
			name:           "Invalid ID",
			id:             "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid ID"}`,
		},
		{
			name:           "Task Not Found",
			id:             "2",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Task not found"}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := http.NewRequest(http.MethodGet, "/tasks/"+tc.id, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Params = gin.Params{{Key: "id", Value: tc.id}}
			planner.GetTask(c)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.JSONEq(t, tc.expectedBody, rr.Body.String())
		})
	}
}

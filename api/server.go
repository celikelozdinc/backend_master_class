package api

import (
	"database/sql"
	"net/http"

	db "code.siemens.com/ozdinc.celikel/backend_master_vlass/internal/db"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our automation tool
type Server struct {
	transaction db.Tx
	router      *gin.Engine
}

type createStudentRequest struct {
	Name   string `json:"name" binding:"required"`
	Number int64  `json:"number" binding:"required"`
	Nation string `json:"nation" binding:"required"`
}

type getStudentRequest struct {
	//since ID is a URI parameter, we cannot get it from the request body
	ID int64 `uri:"id" binding:"required,min=1"` // => smallest possible value of account ID is 1
}

// NewServer creates a new HTTP server and set up routing
func NewServer(tx db.Tx) *Server {
	server := &Server{transaction: tx}
	router := gin.Default()

	// add routes to router
	router.POST("/newStudent", server.createStudent)
	router.GET("/students/:id", server.getStudent)
	router.POST("/grading", server.createGrading)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) createStudent(ctx *gin.Context) {
	var request createStudentRequest
	// For validating the output object to make sure it satisfy the conditions we specified in the binding tag
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateStudentParams{
		Name:   request.Name,
		Number: request.Number,
		Nation: request.Nation,
	}
	account, err := server.transaction.CreateStudent(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) getStudent(ctx *gin.Context) {
	var request getStudentRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	student, err := server.transaction.GetStudent(ctx, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Status Code 404
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		// Status Code 500
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, student)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

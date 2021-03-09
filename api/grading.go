package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "code.siemens.com/ozdinc.celikel/backend_master_vlass/internal/db"
	"github.com/gin-gonic/gin"
)

type gradingRequest struct {
	StudentID int64  `json:"student_id" binding:"required"`
	Grade     int64  `json:"grade" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Number    int64  `json:"number" binding:"required"`
	Nation    string `json:"nation" binding:"required"`
}

func (server *Server) createGrading(ctx *gin.Context) {
	var request gradingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.isStudentValid(ctx, request.StudentID, request.Number) {
		return
	}

	params := db.GradingTxParams{
		StudentID: request.StudentID,
		Grade:     request.Grade,
		Name:      request.Name,
		Number:    request.Number,
		Nation:    request.Nation,
	}

	err := server.transaction.GradingTx(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, `OK`)
}

func (server *Server) isStudentValid(ctx *gin.Context, studentID int64, studentNumber int64) bool {
	student, err := server.transaction.GetStudent(ctx, studentID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if student.Number != studentNumber {
		err := fmt.Errorf("Student [%d] number mismatch: %v != %v", student.ID, student.Number, studentNumber)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}

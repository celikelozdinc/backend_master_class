package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Tx defines methods to be implemented
type Tx interface {
	Querier
	GradingTx(ctx context.Context, args GradingTxParams) error
}

// SQLTx implements a basic database transaction
type SQLTx struct {
	*Queries
	Db *sql.DB
}

type (
	// NewStudent is a reusable function signature
	NewStudent func(*Queries, *GradingTxParams) error

	// NewGrade is a reusable function signature
	NewGrade func(*Queries, *GradingTxParams) error
)

var (
	NS NewStudent
	NG NewGrade
)

func init() {
	NS = func(q *Queries, p *GradingTxParams) error {
		arg := CreateStudentParams{
			Name:   p.Name,
			Nation: p.Nation,
			Number: p.Number,
		}

		_, createErr := q.CreateStudent(context.Background(), arg)
		if createErr != nil {
			return createErr
		}

		return nil
	}

	NG = func(q *Queries, p *GradingTxParams) error {
		arg := CreateGradeParams{
			Grade:     p.Grade,
			StudentID: p.StudentID,
		}

		_, createErr := q.CreateGrade(context.Background(), arg)
		if createErr != nil {
			return createErr
		}

		return nil
	}
}

// NewTx will build a new Tx object
func NewTx(db *sql.DB) Tx {
	return &SQLTx{
		Db:      db,
		Queries: New(db),
	}
}

func (tx *SQLTx) execTx(ctx context.Context, fn func(q *Queries, p *GradingTxParams) error, args *GradingTxParams) error {
	transaction, err := tx.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(transaction) //=> pass a sql.Tx object
	err = fn(q, args)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction err: %v, rb err: %v", err, rollbackErr)
		}
		return err
	}

	return transaction.Commit()
}

// GradingTx inserts both a student and her/his grade
func (tx *SQLTx) GradingTx(ctx context.Context, args GradingTxParams) error {
	if newStudentErr := tx.execTx(ctx, NS, &args); newStudentErr != nil {
		log.Printf("Can not create student, due to %s\n", newStudentErr.Error())
		return newStudentErr
	}

	if newGradeErr := tx.execTx(ctx, NG, &args); newGradeErr != nil {
		log.Printf("Can not create grade, due to %s\n", newGradeErr.Error())
		return newGradeErr
	}

	return nil
}

type GradingTxParams struct {
	StudentID int64  `json:"student_id"`
	Grade     int64  `json:"grade"`
	Name      string `json:"name"`
	Number    int64  `json:"number"`
	Nation    string `json:"nation"`
}

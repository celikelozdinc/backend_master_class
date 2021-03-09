package db

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"testing"

	"code.siemens.com/ozdinc.celikel/backend_master_vlass/util"
	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
)

func Test_CreateStudent(t *testing.T) {
	arg := CreateStudentParams{
		Name:   "Dimitris Diamantidis",
		Nation: "GR",
		Number: 40090,
	}

	account, err := testQueries.CreateStudent(context.Background(), arg)
	if err != nil {
		t.Errorf("Error must be nil")
	}
	t.Log("Student has been created")

	if got, exp := account.Nation, "GR"; got != exp {
		t.Errorf("Unexpected nationality of student. Got :%s, Expected : %s", got, exp)
	}
	t.Logf("Nationality of brand-new student : %s", account.Nation)

	if account.ID == 0 {
		t.Errorf("ID must not be 0")
	}
	t.Logf("Id of brand-new student : %d", account.ID)
}

func Test_GetStudent(t *testing.T) {
	var studentID int64 = 3
	stu, err := testQueries.GetStudent(context.Background(), studentID)

	if err != nil {
		t.Errorf("Error must be nil")
	}

	if got, exp := stu.Nation, "FR"; got != exp {
		t.Errorf("Unexpected nationality of student. Got :%s, Expected : %s", got, exp)
	}

	t.Logf("Name of student : %s", stu.Name)
}

func Test_GradingTransaction(t *testing.T) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, _ := sql.Open(config.DBDriver, config.DBSource)
	defer conn.Close()

	tx := NewTx(conn) // => Transaction object

	errs := make(chan error)

	// Prepare test data
	params := map[int]*GradingTxParams{}

	params[0] = &GradingTxParams{
		StudentID: 1,
		Grade:     0,
		Name:      "Louis Bebe",
		Nation:    "FR",
		Number:    int64(rand.Int()),
	}
	params[1] = &GradingTxParams{
		StudentID: 2,
		Grade:     10,
		Name:      "Louis Bebe",
		Nation:    "FR",
		Number:    int64(rand.Int()),
	}
	params[2] = &GradingTxParams{
		StudentID: 3,
		Grade:     20,
		Name:      "Louis Bebe",
		Nation:    "FR",
		Number:    int64(rand.Int()),
	}

	//run n concurrent grading transaction
	for i := 0; i < 3; i++ {
		go func(counter int) {
			gradingErr := tx.GradingTx(context.Background(), *params[counter])
			errs <- gradingErr
		}(i)
	}

	// check results
	for i := 0; i < 3; i++ {
		err := <-errs
		if err != nil {
			t.Errorf("Error must be nil")
		}

		t.Logf("Checking results for iteration %d : ", i)
		firstGrade, _ := testQueries.GetGradeByStudentID(context.Background(), params[0].StudentID)
		t.Logf("%#v", firstGrade)
		if firstGrade.Grade != params[0].Grade {
			t.Errorf("Expected grade is %d, not %d", 0, firstGrade.Grade)
		}

		secGrade, _ := testQueries.GetGradeByStudentID(context.Background(), params[1].StudentID)
		t.Logf("%#v", secGrade)
		if secGrade.Grade != params[1].Grade {
			t.Errorf("Expected grade is %d, not %d", 10, secGrade.Grade)
		}

		thirdGrade, _ := testQueries.GetGradeByStudentID(context.Background(), params[2].StudentID)
		t.Logf("%#v", thirdGrade)
		if thirdGrade.Grade != params[2].Grade {
			t.Errorf("Expected grade is %d, not %d", 20, thirdGrade.Grade)
		}

	}
}

func init() {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connection, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(connection) //=> pass a sql.DB object

}

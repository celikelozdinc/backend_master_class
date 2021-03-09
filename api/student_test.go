package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	db "code.siemens.com/ozdinc.celikel/backend_master_vlass/internal/db"
	mockdb "code.siemens.com/ozdinc.celikel/backend_master_vlass/internal/db/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_GetStudentAPI(t *testing.T) {
	dummyStudent := db.Student{
		Name:   "MockStudent",
		Nation: "TR",
		Number: 00000,
		ID : 17171717,
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockTx := mockdb.NewMockTx(controller)

	// Define stub
	mockTx.EXPECT().
		GetStudent(gomock.Any(), gomock.Eq(dummyStudent.ID)).
		Times(1).
		Return(dummyStudent, nil)

	server := NewServer(mockTx)
	recorder := httptest.NewRecorder()

	endpoint := fmt.Sprintf("/students/%d", dummyStudent.ID)
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

}

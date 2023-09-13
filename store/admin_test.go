package store

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_Healthz(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	adminStore := NewAdminStore(mockDB)

	t.Run("Should return nil when the query is successful", func(t *testing.T) {

		expectedQuery := "SELECT k8s_pod_id FROM sip_routers"
		rows := sqlmock.NewRows([]string{"k8s_pod_id"})

		mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

		err := adminStore.Healthz()

		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Should return an error when the query fails", func(t *testing.T) {

		expectedQuery := "SELECT k8s_pod_id FROM sip_routers"
		expectedError := sql.ErrNoRows

		mock.ExpectQuery(expectedQuery).WillReturnError(expectedError)

		err := adminStore.Healthz()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})
}

func Test_GetBestRTPProxy(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	adminStore := NewAdminStore(mockDB)

	t.Run("Should return the best RTP proxy host", func(t *testing.T) {

		expectedQuery := "SELECT rtpproxy_sock, set_id, cpu_pct, cpu_used FROM rtpproxy_sockets"
		rows := sqlmock.NewRows([]string{"rtpproxy_sock", "set_id", "cpu_pct", "cpu_used"}).
			AddRow("udp:host1:port1", 1, 10.0, 5.0).
			AddRow("udp:host2:port2", 2, 5.0, 3.0)
		mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

		result, err := adminStore.GetBestRTPProxy()
		assert.NoError(t, err)

		expectedHost := []byte("host2")
		assert.Equal(t, expectedHost, result)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Should return an error when the query fails", func(t *testing.T) {
		expectedQuery := "SELECT rtpproxy_sock, set_id, cpu_pct, cpu_used FROM rtpproxy_sockets"
		expectedError := sql.ErrNoRows

		mock.ExpectQuery(expectedQuery).WillReturnError(expectedError)
		result, err := adminStore.GetBestRTPProxy()

		assert.Error(t, err)
		assert.Nil(t, result)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})
}

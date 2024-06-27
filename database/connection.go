package database

import (
	"database/sql"
	"fmt"
	"time"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sony/gobreaker"
)

// MySQLConn represents a MySQL connection manager with circuit breaker.
type MySQLConn struct {
	db          *sql.DB
	circuit     *gobreaker.CircuitBreaker
	lastAttempt time.Time
}

// NewMySQLConn creates a new MySQL connection manager with circuit breaker using an existing *sql.DB connection.
func NewMySQLConn(db *sql.DB) *MySQLConn {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "MySQL Circuit Breaker",
		MaxRequests: 10,                    // Number of requests allowed before tripping the circuit
		Interval:    600 * time.Second,      // Time period for which to count the number of requests
		Timeout:     10 * time.Second,     // Time period after which to consider a request as timed out
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// If more than 2 requests have failed consecutively, trip the circuit
			return counts.ConsecutiveFailures > 2
		},
	})

	return &MySQLConn{
		db:      db,
		circuit: cb,
	}
}

// Query executes a query against MySQL with circuit breaker protection.
func (m *MySQLConn) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error

	// Execute the query within the circuit breaker
	output, err := m.circuit.Execute(func() (interface{}, error) {
		rows, err = m.db.Query(query, args...)
		return rows, err
	})

	// Update last attempt time
	m.lastAttempt = time.Now()

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return output.(*sql.Rows), nil
}

// QueryRow executes a query that is expected to return at most one row.
func (m *MySQLConn) QueryRow(query string, args ...interface{}) *sql.Row {
	var row *sql.Row

	// Execute the query within the circuit breaker
	_, err := m.circuit.Execute(func() (interface{}, error) {
		row = m.db.QueryRow(query, args...)
		return row, nil
	})

	// Update last attempt time
	m.lastAttempt = time.Now()

	// Return the row and error if any
	if err != nil {
		log.Printf("failed to execute query row: %v", err)
	}
	return row
}

// Prepare prepares a SQL statement for execution.
func (m *MySQLConn) Prepare(query string) (*sql.Stmt, error) {
	var stmt *sql.Stmt
	var err error

	// Prepare the statement within the circuit breaker
	output, err := m.circuit.Execute(func() (interface{}, error) {
		stmt, err = m.db.Prepare(query)
		return stmt, err
	})

	// Update last attempt time
	m.lastAttempt = time.Now()

	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}

	return output.(*sql.Stmt), nil
}

// Exec executes a query that doesn't return rows, such as INSERT, UPDATE, DELETE.
func (m *MySQLConn) Exec(query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	var err error

	// Execute the query within the circuit breaker
	output, err := m.circuit.Execute(func() (interface{}, error) {
		result, err = m.db.Exec(query, args...)
		return result, err
	})

	// Update last attempt time
	m.lastAttempt = time.Now()

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return output.(sql.Result), nil
}

// get the connection
func (m *MySQLConn) GetConnection() (*sql.DB) {
	return m.db;
}

// Close closes the MySQL connection.
func (m *MySQLConn) Close() error {
	return m.db.Close()
}
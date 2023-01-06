package store

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

type DebitStore struct {
	db *sql.DB
}

func NewDebitStore(db *sql.DB) *DebitStore {
	return &DebitStore{
		db: db,
	}
}

func (ds *DebitStore) CreateDebit(rate *utils.CallRate, debit *model.Debit) error {
	minutes := math.Floor(debit.Seconds / 60)
	dollars := minutes * rate.CallRate
	cents := utils.ToCents(dollars)
	now := time.Now()
	stmt, err := ds.db.Prepare("INSERT INTO users_debits (`user_id`, `cents`, `source`, `plan_snapshot`, `module_id`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ?, ? )")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(debit.UserId, cents, debit.Source, debit.PlanSnapshot, debit.ModuleId, now, now)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DebitStore) CreateAPIUsageDebit(workspace *model.Workspace, debitApi *model.DebitAPI) error {
	if debitApi.Type == "TTS" {
		dollars := utils.CalculateTTSCosts(debitApi.Params.Length)
		cents := utils.ToCents(dollars)
		source := fmt.Sprintf("API usage - %s", debitApi.Type)
		now := time.Now()
		stmt, err := ds.db.Prepare("INSERT INTO users_debits (`user_id, `cents`, `source`, `plan_snapshot`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ? )")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(debitApi.UserId, cents, source, workspace.Plan, now, now)
		if err != nil {
			return err
		}
	} else if debitApi.Type == "STT" {
		dollars := utils.CalculateSTTCosts(debitApi.Params.RecordingLength)
		cents := utils.ToCents(dollars)
		source := fmt.Sprintf("API usage - %s", debitApi.Type)
		now := time.Now()
		stmt, err := ds.db.Prepare("INSERT INTO users_debits (`user_id, `cents`, `source`, `plan_snapshot`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ? )")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(debitApi.UserId, cents, source, workspace.Plan, now, now)
		if err != nil {
			return err
		}
	}
	return nil
}

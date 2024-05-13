package store

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

/*
Implementation of Debit Store
*/

type DebitStore struct {
	db *sql.DB
}

func NewDebitStore(db *sql.DB) *DebitStore {
	return &DebitStore{
		db: db,
	}
}

/*
Input: CallRate model, Debit Model
Todo : Create new user_debit and store to db
Output: If success return nil else return err
*/
func (ds *DebitStore) CreateDebit(rate *model.CallRate, debit *model.Debit) error {
	var cents int;

	customizations := utils.GetCustomizationSettings()

	if customizations.BillingFrequency == "PER_MINUTE" {
		minutes := math.Ceil(debit.Seconds / 60)
		dollars := minutes * rate.CallRate
		cents = utils.ToCents(dollars)
	} else if customizations.BillingFrequency == "PER_SECOND" {
		minutes := debit.Seconds / 60
		dollars := minutes * rate.CallRate
		cents = utils.ToCents(dollars)
	}

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

/*
Input: Workspace model, DebitAPI model
Todo : Calculate cents based on debit type and create user_debit
Output: If success return nil else return err
*/
func (ds *DebitStore) CreateAPIUsageDebit(workspace *model.Workspace, debitApi *model.DebitAPI) error {
	// Check DebitType and calcaulte cents individually
	if debitApi.Type == "TTS" {
		dollars := utils.CalculateTTSCosts(debitApi.Params.Length)
		cents := utils.ToCents(dollars)
		source := fmt.Sprintf("API usage - %s", debitApi.Type)
		now := time.Now()
		stmt, err := ds.db.Prepare("INSERT INTO users_debits (`user_id`, `cents`, `source`, `plan_snapshot`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ? )")
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
		stmt, err := ds.db.Prepare("INSERT INTO users_debits (`user_id`, `cents`, `source`, `plan_snapshot`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ? )")
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

package customizations

import (
	"github.com/sirupsen/logrus"
	"lineblocs.com/api/utils"
	helpers "github.com/Lineblocs/go-helpers"
)


func GetCustomizationSettings() (*helpers.CustomizationSettings, error) {
	db, err := helpers.CreateDBConn()
	if err != nil {
		utils.Log(logrus.PanicLevel, err.Error())
		panic(err)
	}

	row := db.QueryRow("SELECT invoice_due_date_enabled, COALESCE(0,invoice_due_num_days), billing_frequency, customer_satisfaction_survey_enabled, customer_satisfaction_survey_url FROM customizations")
	if err != nil {
		return nil, err
	}

	value := helpers.CustomizationSettings{}
	err = row.Scan(&value.InvoiceDueDateEnabled, 
		&value.InvoiceDueNumDays,
		&value.BillingFrequency,
		&value.CustomerSatisfactionSurveyEnabled,
		&value.CustomerSatisfactionSurveyUrl,
	)
	if err != nil {
		utils.Log(logrus.PanicLevel, err.Error())
		return nil, err
	}

	return &value,nil
}

// loadDataFromDatabase loads data from the SQL database
func loadDataFromDatabase() (*helpers.CustomizationSettings, error) {
	record, err := GetCustomizationSettings()
	if err != nil {
		return nil,err
	}

	return record,nil
}

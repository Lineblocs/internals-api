package customizations

import (
	"fmt"
	"sync"
	helpers "github.com/Lineblocs/go-helpers"
)

var (
	instance    *helpers.CustomizationSettings
	once        sync.Once
	initialized bool
)

// GetInstance returns the singleton instance of CustomizationSettings
// TODO: fix this so it only sets the instance variable one time.
func GetInstance() (*helpers.CustomizationSettings, error) {
	var err error

	data, err := loadDataFromDatabase()
	if err != nil {
		err = fmt.Errorf("failed to load data: %w", err)
		return nil, err
	}

	initialized = true
	return data, nil
}

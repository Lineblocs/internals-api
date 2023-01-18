package carrier

import "lineblocs.com/api/helpers"

/*
Interface of Carrier Store.
Implementation of Carrier Store is located /store/carrier
*/
type Store interface {
	CreateSIPReport(string, string) error
	CreateRoutingFlow(*string, *string, *string) (*helpers.Flow, error)
	StartProcessingFlow(*helpers.Flow, map[string]string) ([]*helpers.RoutablePSTNProvider, error)
}

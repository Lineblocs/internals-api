package admin

/*
Interface of Admin Store.
Implementation of Admin Store is located /store/admin
*/
type AdminStoreInterface interface {
	GetBestRTPProxy() ([]byte, error)
	Healthz() error
}

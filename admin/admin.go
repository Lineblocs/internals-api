package admin

/*
Interface of Admin Store.
Implementation of Admin Store is located /store/admin
*/
type Store interface {
	GetBestRTPProxy() ([]byte, error)
	Healthz() error
}

package admin

type Store interface {
	GetBestRTPProxy() ([]byte, error)
}

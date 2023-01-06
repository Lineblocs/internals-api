package store

import "database/sql"

type AdminStore struct {
	db *sql.DB
}

func NewAdminStore(db *sql.DB) *AdminStore {
	return &AdminStore{
		db: db,
	}
}

func (as *AdminStore) GetBestRTPProxy() ([]byte, error) {
	results, err := as.db.Query(`SELECT rtpproxy_sock, set_id, cpu_pct, cpu_used FROM rtpproxy_sockets`)
	// Execute the query
	if err != nil {
		return nil, err
	}
	defer results.Close()
	var lowestCPU *float64 = nil
	var socketAddrResult string
	for results.Next() {
		var rtpSock string
		var setId int
		var cpuPct float64
		var cpuUsed float64
		err = results.Scan(&rtpSock, &setId, &cpuPct, &cpuUsed)
		if err != nil {
			return nil, err
		}
		if lowestCPU == nil || cpuPct < *lowestCPU {
			lowestCPU = &cpuPct
			socketAddrResult = rtpSock
		}
	}
	return []byte(socketAddrResult), nil
}

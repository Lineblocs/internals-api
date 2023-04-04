package store

import (
	"database/sql"
	"strings"
)

/*
Implementation of Admin Store
*/

type AdminStore struct {
	db *sql.DB
}

func NewAdminStore(db *sql.DB) *AdminStore {
	return &AdminStore{
		db: db,
	}
}

/*
Todo: Check SIP Routers is Healthy
Output: If excute query success return nil else return err
*/
func (as *AdminStore) Healthz() error {
	// Execute the query...
	results, err := as.db.Query("SELECT k8s_pod_id FROM sip_routers")
	if err != nil {
		return err
	}
	defer results.Close()
	return nil
}

/*
Input: EmailInfo model
Todo : Send Email to Lineblocs Contact
Output: First value: rtpSock, Second value: err
If success return (rtpSock,nil) else return (nil, err)
*/
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
	// socket looks like
	// udp:host:port
	// we only need the host for OpenSIPs
	splitted := strings.Split(socketAddrResult, ":")
	socketHost := splitted[1]
	return []byte(socketHost), nil
}

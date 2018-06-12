package info

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocql/gocql"
)

const keyspaceName = "myinfos"

var cluster *gocql.ClusterConfig

func init() {
	var err error
	err = createKeyspace()
	if err == nil {
		cluster = gocql.NewCluster("127.0.0.1")
		cluster.Port = 9042
		cluster.ProtoVersion = 4
		cluster.Keyspace = keyspaceName

		err = checkDatabase()
	}
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// GetMachines return all machines in the database.
func GetMachines() ([]Machine, error) {
	session, err := cluster.CreateSession()
	defer session.Close()
	if err != nil {
		return nil, err
	}

	machines := make([]Machine, 0)

	var ipAddress, hostname string
	var uptime, runningSince int64
	iter := session.Query("SELECT ip_address, hostname, uptime, running_since FROM myinfos.host").Iter()
	for iter.Scan(&ipAddress, &hostname, &uptime, &runningSince) {
		machine := Machine{
			IPAddress:    ipAddress,
			Hostname:     hostname,
			Uptime:       uptime,
			RunningSince: runningSince,
		}
		machines = append(machines, machine)
	}
	return machines, nil
}

// SaveMachine try to save the strcut
func SaveMachine(machine Machine) error {
	session, err := cluster.CreateSession()
	defer session.Close()
	if err != nil {
		return err
	}

	json, errJ := json.Marshal(machine)
	if errJ != nil {
		return errJ
	}

	err = session.Query("INSERT INTO host JSON ?", string(json)).Exec()
	return err
}

func createKeyspace() error {
	clusterTemp := gocql.NewCluster("127.0.0.1")
	clusterTemp.Port = 9042
	clusterTemp.ProtoVersion = 4
	clusterTemp.Keyspace = "system"

	session, err := clusterTemp.CreateSession()
	defer session.Close()
	if err != nil {
		return err
	}

	count := 0
	session.Query("SELECT count(1) as count FROM system_schema.keyspaces WHERE  keyspace_name = ? LIMIT 1", keyspaceName).Scan(&count)
	if count == 0 {
		cmd := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = {'class' : 'SimpleStrategy','replication_factor' : 1}", keyspaceName)
		err = session.Query(cmd).Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func checkDatabase() error {
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var metadata *gocql.KeyspaceMetadata
	metadata, err = session.KeyspaceMetadata(keyspaceName)
	if err != nil {
		return err
	}

	if _, exists := metadata.Tables["host"]; !exists {
		cmd := fmt.Sprintf(`CREATE TABLE %s.host (
				ip_address text PRIMARY KEY,
				hostname text,
				uptime varint,
				running_since varint,
			) WITH comment='List of hosts with the uptime machine and uptime of application'`, keyspaceName)
		if err = session.Query(cmd).Exec(); err != nil {
			return err
		}
	}

	return nil
}

package info

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocql/gocql"
)

const keyspaceName = "myinfos"

var (
	cluster *gocql.ClusterConfig
	host    string
	port    int
)

// InitDatabase the access of database.
func InitDatabase(dbhost string, dbport int) {
	log.Printf("Connecting to Cassandra database in: %s:%d\n", dbhost, dbport)
	host = dbhost
	port = dbport

	var err error
	err = createKeyspace()
	if err == nil {
		cluster = gocql.NewCluster(host)
		cluster.Port = port
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

	log.Println("Getting a list of machines...")

	machines := make([]Machine, 0)

	var ipAddress, hostname string
	var uptime, lastPing int64
	iter := session.Query("SELECT ip_address, hostname, uptime, last_ping FROM myinfos.host").Iter()
	for iter.Scan(&ipAddress, &hostname, &uptime, &lastPing) {
		machine := Machine{
			IPAddress: ipAddress,
			Hostname:  hostname,
			Uptime:    uptime,
			LastPing:  lastPing,
		}
		machines = append(machines, machine)
	}

	log.Printf("Machines:\n%v\n", machines)
	log.Println("Getting a list of machines...Okay")

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

	log.Printf("Saving the machine info:%v\n", machine)
	err = session.Query("INSERT INTO host JSON ?", string(json)).Exec()
	return err
}

func createKeyspace() error {
	clusterTemp := gocql.NewCluster(host)
	clusterTemp.Port = port
	clusterTemp.ProtoVersion = 4
	clusterTemp.Keyspace = "system"

	session, err := clusterTemp.CreateSession()
	defer session.Close()
	if err != nil {
		return err
	}

	log.Println("Checking the keyspace...")

	count := 0
	session.Query("SELECT count(1) as count FROM system_schema.keyspaces WHERE  keyspace_name = ? LIMIT 1", keyspaceName).Scan(&count)
	if count == 0 {
		log.Printf("    Looks like first access, creating the keyspace '%s'\n", keyspaceName)
		cmd := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = {'class' : 'SimpleStrategy','replication_factor' : 1}", keyspaceName)
		err = session.Query(cmd).Exec()
		if err != nil {
			return err
		}
	}
	log.Println("Checking the keyspace...Okay")
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

	log.Println("Checking the tables...")

	if _, exists := metadata.Tables["host"]; !exists {
		log.Println("    Looks like the table 'host' not exist yet. Creating the table...")
		cmd := fmt.Sprintf(`CREATE TABLE %s.host (
				ip_address text PRIMARY KEY,
				hostname text,
				uptime varint,
				last_ping varint,
			) WITH comment='List of hosts with the uptime machine and uptime of application'`, keyspaceName)
		if err = session.Query(cmd).Exec(); err != nil {
			return err
		}
	}

	log.Println("Checking the tables...Okay")
	return nil
}

package bigdb

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/rs/zerolog/log"
)

//CassandraConfig is used to init connection to cassandra instance.
type CassandraConfig struct {
	IP         []string
	cqlVersion string
	port       int
	user       string
	pass       string
}

//DefaultClusterConfig can be used to init CassandraConfig with default config.
func DefaultClusterConfig() *CassandraConfig {
	return &CassandraConfig{
		IP:         []string{"127.0.0.1"},
		cqlVersion: "4.0.0",
		port:       9042,
		user:       "cassandra",
		pass:       "cassandra",
	}
}

//NewCassandraSession initalizes cassandra
func NewCassandraSession(config *CassandraConfig) *gocql.Session {
	cluster := gocql.NewCluster(config.IP...)
	cluster.CQLVersion = config.cqlVersion
	cluster.Port = config.port
	cluster.Timeout = time.Second * 5
	cluster.ConnectTimeout = time.Second * 5
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.user,
		Password: config.pass,
	}
	session, err := cluster.CreateSession()
	if err != nil {
		log.Error().Msgf("error in cassandra session create: %v", err)
	}
	return session
}

package utils

import (
	"github.com/gocql/gocql"
	"github.com/rs/zerolog/log"
)

func initCassandra(cassPort int) *gocql.Session {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.CQLVersion = "4.0.0"
	cluster.Port = cassPort
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "cassandra",
	}
	session, err := cluster.CreateSession()
	if err != nil {
		log.Error().Msgf("error in cassandra session create: %v", err)
	}
	return session
}

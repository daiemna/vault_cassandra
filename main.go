package main

import (
	"fmt"
	"time"

	"github.com/daiemna/vault_cassandra/internal/bigdb"
	"github.com/daiemna/vault_cassandra/internal/utils"
	"github.com/daiemna/vault_cassandra/internal/vault"
	"github.com/gocql/gocql"
)

var log = utils.SetupLogger(true)

func doseRoleExist(roleName string, session *gocql.Session) error {
	var roleID string
	err := session.Query("SELECT role FROM system_auth.roles WHERE role = ?", roleName).Scan(&roleID)
	if err != nil {
		return fmt.Errorf("error in role query: %v", err)
	}
	log.Debug().Msgf("roleName: %s", roleName)
	if roleID != roleName {
		return fmt.Errorf("role not present")
	}
	return nil
}

func testCassandraConnection(session *gocql.Session) {
	ksMeta, err := session.KeyspaceMetadata("system_auth")
	if err != nil {
		log.Fatal().Msgf("error in meta data retrival: %v", err)
	}
	log.Debug().Msgf("ks StrategyClass : %s", ksMeta.StrategyClass)
}

func testVaultCassandra(session *gocql.Session, vaultClient *vault.Client) {
	const delay = 1
	defer log.Info().Msg("Test passed!")
	secret, err := vaultClient.Read("database/creds/my-role")
	if err != nil {
		log.Debug().Msgf("Error in vault read : %v", err)
	}
	log.Debug().Msgf("secrets : %v", secret)
	leaseDuration := time.Duration(secret.LeaseDuration+delay) * time.Second
	roleName := secret.Data["username"].(string)

	if err := doseRoleExist(roleName, session); err != nil {
		log.Fatal().Msgf("Error in role check, before role expiry: %v", err)
	}

	log.Info().Msg("Waiting for credentials to expire ... bis bald!")
	time.Sleep(leaseDuration)

	if err := doseRoleExist(roleName, session); err != nil {
		log.Info().Msgf("Error in role check, after role expiry: %v", err)
	}
}

func main() {
	cassConf := bigdb.DefaultClusterConfig()
	cassSession := bigdb.NewCassandraSession(cassConf)
	log.Debug().Msg("cassandra init done")

	vaultClient, err := vault.NewRootClient()

	if err != nil {
		log.Fatal().Msgf("error creating vault client: %v", err)
	}

	testCassandraConnection(cassSession)

	testVaultCassandra(cassSession, vaultClient)
}

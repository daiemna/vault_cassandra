package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
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

func testCassandraCredentials(session *gocql.Session, vaultClient *vault.Client) {
	const delay = 1

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
		if strings.Contains(err.Error(), "not found") {
			log.Info().Msg("Role not found, as expected.")
			log.Info().Msg("Test Done!")
			return
		}
		log.Fatal().Msgf("Error in role check, after role expiry: %v", err)
	}
}

func testCassandraCredentialsMultiRead(session *gocql.Session, vaultClient *vault.Client, readCount int) error {
	const delay = 1

	secret, err := vaultClient.Read("database/creds/my-role")
	if err != nil {
		return fmt.Errorf("Error in vault read : %v", err)
	}
	roleName := secret.Data["username"].(string)

	for i := 0; i < readCount; i++ {
		err := doseRoleExist(roleName, session)
		if err != nil {
			return fmt.Errorf("Error in role read check: %v", err)
		}
	}
	return nil
}

func testVaultCassandra(routineNumber, readCount int, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	cassConf := bigdb.DefaultClusterConfig()
	cassSession := bigdb.NewCassandraSession(cassConf)
	// log.Debug().Msgf("cassandra init done routine num %d", routineNumber)

	vaultClient, err := vault.NewRootClient()

	if err != nil {
		log.Fatal().Msgf("error creating vault client: %v", err)
	}

	err = testCassandraCredentialsMultiRead(cassSession, vaultClient, readCount)
	if err != nil {
		log.Fatal().Msgf("error testing cassandra: %v", err)
	}

}

func main() {
	routineCount := 50
	readCount := 50
	var waitGroup sync.WaitGroup
	os.Setenv("VAULT_TOKEN", "testtoken")
	for i := 0; i < routineCount; i++ {
		go testVaultCassandra(i, readCount, &waitGroup)
	}
	waitGroup.Wait()
}

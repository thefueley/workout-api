package database

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/rs/zerolog/log"
)

// NewDatabase : return pointer to db object
func NewDatabase() (*aztables.ServiceClient, error) {
	log.Info().Msg("Setting up DB connection")

	connectionString := os.Getenv("AZ_TABLE_CONN_STR")

	service, err := aztables.NewServiceClientFromConnectionString(connectionString, nil)

	if err != nil {
		log.Error().Msg(err.Error())
	}
	return service, nil
}

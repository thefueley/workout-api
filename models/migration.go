package models

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/rs/zerolog/log"
)

// MigrateDb : migrate db and create comment table
func MigrateDB(svc *aztables.ServiceClient) error {
	tableName := os.Getenv("AZ_TABLE_NAME")
	// check if az table exists
	filter := fmt.Sprintf("TableName eq '%s'", tableName)
	options := &aztables.ListTablesOptions{
		Filter: &filter,
	}

	var tableExists aztables.ListTablesPage
	pager := svc.ListTables(options)

	for pager.NextPage(context.TODO()) {
		tableExists = pager.PageResponse()
		fmt.Printf("Received: %+v\n", len(tableExists.Tables))
	}

	if len(tableExists.Tables) > 0 {
		log.Info().Msg("Table already exists.")
	} else {
		// create table
		_, err := svc.CreateTable(context.TODO(), tableName, nil)
		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}
	}

	return nil
}

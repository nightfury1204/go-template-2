package main

import (
	"context"
	"log"

	"bitbucket.org/evaly/go-boilerplate/config"
	infraMongo "bitbucket.org/evaly/go-boilerplate/infra/mongo"
	"bitbucket.org/evaly/go-boilerplate/repo"

	"github.com/spf13/cobra"
)

var repos []repo.Repo

var migrationRoot = &cobra.Command{
	Use:   "migration",
	Short: "Run database migrations",
	Long:  `Migration is a tool to generate and modify databse tables`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfgMongo := config.GetMongo(cfgPath)
		cfgDBTable := config.GetTable(cfgPath)

		ctx := context.Background()

		//lgr := logger.DefaultOutStructLogger

		db, err := infraMongo.New(ctx, cfgMongo.URL, cfgMongo.DBName, cfgMongo.DBTimeOut)
		if err != nil {
			return err
		}
		defer db.Close(ctx)

		brandRepo := repo.NewBrand(cfgDBTable.BrandCollectionName, db)

		repos = []repo.Repo{
			brandRepo,
		}

		return nil
	}}

func init() {
	migrationRoot.PersistentFlags().StringVarP(&cfgPath, "config", "c", "config.yaml", "config file path")
}

var migrationUp = &cobra.Command{
	Use:   "up",
	Short: "Populate tables in database",
	Long:  `Populate tables in database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Populating database indices...")
		for _, t := range repos {
			if err := t.EnsureIndices(); err != nil {
				log.Println(err)
			}
		}
		log.Println("Populating database indices successfully...")
		return nil
	},
}

var migrationDown = &cobra.Command{
	Use:   "down",
	Short: "Drop tables from database",
	Long:  `Drop tables from database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Dropping database table...")
		for _, t := range repos {
			if err := t.DropIndices(); err != nil {
				log.Println(err)
			}
		}

		log.Println("Database dopped successfully!")
		return nil
	},
}

func init() {
	migrationRoot.AddCommand(
		migrationUp,
		migrationDown,
	)
}

/*
Copyright © 2024 João Guilhermer joaogbsczip@gmail.com
*/
package cmd

import (
	"database/sql"
	"fmt"

	"github.com/joaosczip/timescale/configs"
	"github.com/joaosczip/timescale/internal/commands"
	"github.com/joaosczip/timescale/internal/common"
	"github.com/joaosczip/timescale/internal/database/repositories"
	"github.com/joaosczip/timescale/internal/handlers"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

// queryMetricsCmd represents the queryMetrics command
var queryMetricsCmd = &cobra.Command{
	Use:   "queryMetrics",
	Short: "Issue queries from the input file against the database to generate some metrics",
	Long: `This command will use the input file to issue queries against the database.
The queries will be used to generate some metrics that will evaluate the performance from both the database and the application.
`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := configs.Load(".")

		if err != nil {
			panic(err)
		}

		db, err := sql.Open(
			"postgres",
			fmt.Sprintf(
				"user=%s dbname=%s password=%s host=%s port=%d sslmode=%s",
				config.Database.User, config.Database.Name, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.SSLModel,
			),
		)

		if err != nil {
			panic(err)
		}

		defer db.Close()

		db.SetMaxOpenConns(config.Database.MaxOpenConns)

		err = cmd.ParseFlags(args)

		if err != nil {
			panic(err)
		}

		queryParamsPath, err := cmd.Flags().GetString("query-params-path")

		if err != nil {
			panic(err)
		}

		repository := repositories.NewCpuUsageRepository(db)
		handler := handlers.NewCollectCpuUsageMetricsHandler(repository)
		csvReader := common.NewQueryParamsCsvReader()
		stdoutWriter := common.NewTableWriter()

		command := commands.NewDatabaseMetricsCommand(handler, csvReader, stdoutWriter)
		command.Run(queryParamsPath)
	},
}

func init() {
	rootCmd.AddCommand(queryMetricsCmd)

	queryMetricsCmd.Flags().StringP("query-params-path", "p", "", "path to the file containing the query params to be issued against the database")
}

/*
Copyright © 2024 João Guilhermer joaogbsczip@gmail.com
*/
package cmd

import (
	"fmt"

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
		fmt.Println("queryMetrics called")
	},
}

func init() {
	rootCmd.AddCommand(queryMetricsCmd)

	queryMetricsCmd.Flags().StringP("query-params-path", "p", "", "path to the file containing the query params to be issued against the database")
}

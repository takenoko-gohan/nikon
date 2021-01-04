/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/takenoko-gohan/nikon/internal/indices"
)

// indicesCmd represents the indices command
var indicesCmd = &cobra.Command{
	Use:   "indices <http://localhost:9200>",
	Short: "Command to get index list of target Elasticsearch",
	Long: `Command to get index list of target Elasticsearch.
If you do not specify the target,
it will try to get the index list from "http: // localhost: 9200".`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			indices.GetIndexList(args[0])
		} else {
			indices.GetIndexList("http://localhost:9200")
		}
	},
}

func init() {
	rootCmd.AddCommand(indicesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// indicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// indicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

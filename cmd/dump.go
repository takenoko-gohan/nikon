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
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/takenoko-gohan/nikon/internal/dump"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump <target index name>",
	Short: "Command to save the target index to a file.",
	Long: `Command to save the target index to a file.
The document is saved in a file in NDJSON structure.`,

	Run: func(cmd *cobra.Command, args []string) {
		h, err := cmd.Flags().GetString("host")
		if err != nil {
			log.Fatal(err)
		}
		s, err := cmd.Flags().GetInt("size")
		if err != nil {
			log.Fatal(err)
		}
		if len(args) > 0 {
			dump.SavingIndex(h, args[0], s)
		} else {
			fmt.Println("Please specify the target index.")
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dumpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dumpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	dumpCmd.Flags().StringP("host", "h", "http://localhost:9200", "Specify the target Elasticsearch")
	dumpCmd.Flags().IntP("size", "s", 100, "Specify the number of items to be acquired at one time")
}

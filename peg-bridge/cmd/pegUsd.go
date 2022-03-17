/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"net/http"

	"github.com/linkpoolio/bridges"
	"github.com/spf13/cobra"
)

// pegUsdCmd represents the pegUsd command
var pegUsdCmd = &cobra.Command{
	Use:   "pegUsd",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		pegUSD(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(pegUsdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pegUsdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pegUsdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type PegUSDBridge struct{}

func (gs *PegUSDBridge) Run(h *bridges.Helper) (interface{}, error) {
	obj := make(map[string]interface{})
	err := h.HTTPCall(
		http.MethodGet,
		"https://pw.gton.capital/rpc/gc-current-peg-usd",
		&obj,
	)
	return obj, err
}

func (gs *PegUSDBridge) Opts() *bridges.Opts {
	return &bridges.Opts{
		Name:   "PegBase",
		Lambda: true,
	}
}

func pegUSD(cmd *cobra.Command, args []string) {
	bridges.NewServer(&PegUSDBridge{}).Start(8080)
}

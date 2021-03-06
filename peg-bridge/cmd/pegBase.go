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
	"fmt"
	"net/http"

	"github.com/linkpoolio/bridges"
	"github.com/spf13/cobra"
)

// pegBaseCmd represents the pegBase command
var pegBaseCmd = &cobra.Command{
	Use:   "pegBase",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		pegBase(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(pegBaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pegBaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pegBaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type PegBaseBridge struct{}

func (gs *PegBaseBridge) Run(h *bridges.Helper) (interface{}, error) {
	/*
		Method base-price (https://pw-rs.gton.capital/rpc/base-price)
		Method owned/base-pool-lps (https://pw-rs.gton.capital/rpc/owned/base-pool-lps)
		Method owned/usd-pool-lps (https://pw-rs.gton.capital/rpc/owned/usd-pool-lps)
		Method base-liquidity (https://pw-rs.gton.capital/rpc/base-liquidity)
		Method usd-liquidity (https://pw-rs.gton.capital/rpc/usd-liquidity)
		Method base-pool-lps (https://pw-rs.gton.capital/rpc/base-pool-lps)
		Method usd-pool-lps (https://pw-rs.gton.capital/rpc/usd-pool-lps)
		Method gc-pol (https://pw-rs.gton.capital/rpc/gc-pol)
	*/
	from := h.GetParam("from")
	url := fmt.Sprintf("https://pw-rs.gton.capital/rpc/%s", from)

	obj := make(map[string]interface{})
	err := h.HTTPCall(
		http.MethodGet,
		url,
		&obj,
	)
	return obj, err
}

func (gs *PegBaseBridge) Opts() *bridges.Opts {
	return &bridges.Opts{
		Name:   "PegBase",
		Lambda: true,
	}
}

func pegBase(cmd *cobra.Command, args []string) {
	bridges.NewServer(&PegBaseBridge{}).Start(8080)
}

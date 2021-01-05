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
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"myctl/services"
	"strings"
)

// applyCmd represents the apply command
var f string
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply a file",
	Long:  `apply a file`,
	Run: func(cmd *cobra.Command, args []string) {
		if strings.TrimSpace(f) != "" {
			file, err := ioutil.ReadFile(f)
			if err != nil {
				fmt.Println("file not exists")
				return
			}

			var users services.UsersInputRequest
			err = yaml.Unmarshal(file, &users)
			fmt.Println(users.Users)
			result, err := userClient.AddUsers(context.Background(), &users)
			if result.Success {
				fmt.Println("update success")
			} else {
				fmt.Println("update fail")
			}
		} else {
			fmt.Println("file name is error")
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().StringVarP(&f, "file", "f", "", "input file")
}

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

			data := make(map[interface{}]interface{})
			err = yaml.Unmarshal(file, &data)
			if err != nil {
				fmt.Println("unmarshal file error")
			}
			fmt.Println(data)
			var users []*services.UserInputRequest
			us := data["users"].([]interface{})

			for i := 0; i < len(us); i++ {
				var id int32
				m := us[i].(map[interface{}]interface{})
				if m["id"] != nil {
					id = int32(m["id"].(int))
				}
				users = append(users, &services.UserInputRequest{
					Username: m["username"].(string),
					Password: m["password"].(string),
					Tel:      m["tel"].(string),
					Email:    m["email"].(string),
					Id:       id,
				})
			}
			result, err := userClient.AddUsers(context.Background(), &services.UsersInputRequest{Users: users})
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

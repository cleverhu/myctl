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
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"log"
	"myctl/services"
)

var format string
var page int32
var size int32
var userClient services.UserServiceClient
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get users",
	Long:  `get users`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if size <= 1 {
			size = 1
		}
		if page <= 1 {
			page = 1
		}
		if len(args) == 1 && args[0] == "users" {
			users, _ := userClient.GetUserList(context.Background(), &services.UserRequest{
				Page:   page,
				Size:   size,
				Search: "",
			})
			if format == "json" {
				data, _ := json.Marshal(users)
				fmt.Printf("%s\n", string(data))
			} else if format=="yaml"{
				data, _ := yaml.Marshal(users)
				fmt.Printf("%s\n", string(data))
			} else {
				fmt.Printf("%-10s\t%-10s%-10s\t%-10s\t%-20s\t%-10s\t\n", "ID", "Name", "Password", "Tel", "Email", "Time")
				for i := 0; i < len(users.Users); i++ {
					m := users.Users[i]
					fmt.Printf("%-10v\t%-10v%-10v\t%-10v\t%-20v\t%-10v\t\n", m.Id, m.Username, m.Password, m.Tel, m.Email, m.CreateTime)
				}
			}
		}

	},
}

func init() {
	conn, err := grpc.Dial("101.132.107.3:8088", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	//defer conn.Close()

	userClient = services.NewUserServiceClient(conn)

	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&format, "format", "o", "", "fmt users format")
	getCmd.Flags().Int32VarP(&size, "size", "s", 10, "query size")
	getCmd.Flags().Int32VarP(&page, "page", "p", 1, "query page")
}

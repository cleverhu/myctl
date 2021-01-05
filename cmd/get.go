package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"myctl/services"
	"os"
	"time"
)

type Server struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
}

type T struct {
	Server  []Server `yaml:"server"`
	Current int      `yaml:"current"`
}

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
			fmt.Println(users)
			if format == "json" {
				data, _ := json.MarshalIndent(users, "", "    ")
				fmt.Printf("%s\n", string(data))
			} else if format == "yaml" {
				data, _ := yaml.Marshal(users)
				fmt.Printf("%s\n", string(data))
			} else {
				fmt.Printf("%-10s\t%-20s%-10s\t%-10s\t%-20s\t%-10s\t\n", "ID", "Name", "Password", "Tel", "Email", "Time")
				for i := 0; i < len(users.Users); i++ {
					m := users.Users[i]
					fmt.Printf("%-10v\t%-20v%-10v\t%-10v\t%-20v\t%-10v\t\n", m.Id, m.Username, m.Password, m.Tel, m.Email, m.CreateTime)
				}
			}
		}

	},
}

func init() {
	d := &T{}
	file, _ := ioutil.ReadFile("config.yaml")
	err := yaml.Unmarshal(file, &d)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().Unix())
	if d.Current > len(d.Server)-1 {
		d.Current = 0
	}
	add := d.Server[d.Current].Address
	//fmt.Println(add)
	conn, err := grpc.Dial(add, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	d.Current = d.Current + 1
	data, _ := yaml.Marshal(d)
	ioutil.WriteFile("config.yaml", data, os.ModePerm)
	//fmt.Println(conn)
	//defer conn.Close()

	userClient = services.NewUserServiceClient(conn)

	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&format, "format", "o", "", "fmt users format")
	getCmd.Flags().Int32VarP(&size, "size", "s", 10, "query size")
	getCmd.Flags().Int32VarP(&page, "page", "p", 1, "query page")
}

package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/hashicorp/consul/api"
)

type runStatus int

const (
	runOK runStatus = iota
	runErr
	runExit
)

func help() {
	fmt.Println("list, get, keys, list, deltree")
}

func client() *api.Client {
	client, _ := api.NewClient(api.DefaultConfig())
	return client
}

func kv() *api.KV {
	return client().KV()
}

func agent() *api.Agent {
	return client().Agent()
}

func kvkeys(c *cli.Context) {
	prefix := c.Args().First()
	keys, _, _ := kv().Keys(prefix, "", nil)

	for _, key := range keys {
		fmt.Printf("%s\n", key)
	}
}

func kvDelTree(c *cli.Context) {
	prefix := c.Args().First()

	if prefix == "" {
		fmt.Print("A key prefix must be specified\n")
		return
	}

	_, err := kv().DeleteTree(prefix, nil)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}

	fmt.Printf("Deleted %s\n", prefix)
}

func kvlist(c *cli.Context) {
	prefix := c.Args().First()
	pairs, _, _ := kv().List(prefix, nil)

	for _, pair := range pairs {
		fmt.Printf("Key: %s, Value: %s\n", pair.Key, string(pair.Value))
	}
}

func kvget(c *cli.Context) {
	key := c.Args().First()

	pair, _, err := kv().Get(key, nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	if pair == nil {
		fmt.Printf("Unable to find key '%s'\n", key)
		return
	}

	fmt.Printf("%s", string(pair.Value))
}

func kvset(c *cli.Context) {
	key := c.Args().Get(0)
	val := c.Args().Get(1)

	if key == "" {
		fmt.Print("A key must be specified\n")
		return
	}

	p := &api.KVPair{Key: key, Value: []byte(val)}
	_, err := kv().Put(p, nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	fmt.Printf("Key: %s, Value: %s\n", key, val)
}

func main() {

	app := cli.NewApp()
	app.Name = "consul-kv"
	app.Usage = "consul kv cli client"
	app.Version = "0.1.0"
	app.Author = "Acaleph"
	app.Email = ""

	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{
	// 		Name:  "node",
	// 		Value: "127.0.0.1:8500",
	// 		Usage: "address:port of the consul node",
	// 	},
	// }

	app.Commands = []cli.Command{
		{
			Name:   "set",
			Usage:  "Set a key & value in the kv store",
			Action: kvset,
		},
		{
			Name:   "get",
			Usage:  "Get a value from the kv store",
			Action: kvget,
		},
		{
			Name:   "keys",
			Usage:  "List keys in the kv store",
			Action: kvkeys,
		},
		{
			Name:   "list",
			Usage:  "List keys and values in the kv store",
			Action: kvlist,
		},
		{
			Name:   "deltree",
			Usage:  "Delete trees in the kv store",
			Action: kvDelTree,
		},
	}

	app.Run(os.Args)
}

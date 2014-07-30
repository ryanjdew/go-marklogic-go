package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	goMarkLogicGo "github.com/ryanjdew/go-marklogic-go"
)

var host string
var port int64
var username string
var password string
var auth string
var queryStr string

func main() {
	flag.StringVar(&host, "host", "localhost", "MarkLogic REST Host")
	flag.Int64Var(&port, "port", 8041, "MarkLogic REST Port")
	flag.StringVar(&username, "username", "admin", "MarkLogic REST Username")
	flag.StringVar(&password, "password", "admin", "MarkLogic REST Password")
	flag.StringVar(&auth, "auth", "basic", "MarkLogic REST Authentication method")
	flag.StringVar(&queryStr, "query", "query", "Search query file")
	flag.Parse()
	var authType int
	if auth == "basic" {
		authType = goMarkLogicGo.BasicAuth
	} else if auth == "digest" {
		authType = goMarkLogicGo.DigestAuth
	} else {
		authType = goMarkLogicGo.None
	}
	client, err := goMarkLogicGo.NewClient(host, port, username, password, authType)
	fmt.Print("Client:\n")
	fmt.Print(spew.Sdump(client))
	if err != nil {
		log.Fatal(err)
	}
	query := goMarkLogicGo.NewQuery(goMarkLogicGo.XML)
	query.Queries = []interface{}{
		&goMarkLogicGo.TermQuery{
			Terms: []string{queryStr},
		},
	}
	fmt.Print("decoded query:\n")
	fmt.Print(spew.Sdump(query))
	resp, err := client.StructuredSearch(query, 1, 10)
	fmt.Print("decoded response:\n")
	fmt.Print(spew.Sdump(resp))
}

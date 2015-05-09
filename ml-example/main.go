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
	flag.Int64Var(&port, "port", 8050, "MarkLogic REST Port")
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
	query := goMarkLogicGo.Query{Format: goMarkLogicGo.XML}
	query.Queries = []interface{}{
		goMarkLogicGo.TermQuery{
			Terms: []string{queryStr},
		},
	}

	qh := goMarkLogicGo.QueryHandle{}
	qh.Decode(query)
	fmt.Print("decoded query:\n")
	fmt.Print(spew.Sdump(qh.Serialized()))
	respHandle := goMarkLogicGo.ResponseHandle{}
	err = client.StructuredSearch(&qh, 1, 10, &respHandle)
	resp := respHandle.Get()
	fmt.Print("decoded response:\n")
	fmt.Print(spew.Sdump(resp))
	sugRespHandle := goMarkLogicGo.SuggestionsResponseHandle{}
	err = client.StructuredSuggestions(&qh, queryStr, 10, "", &sugRespHandle)
	sugResp := sugRespHandle.Serialized()
	fmt.Print("decoded response:\n")
	fmt.Print(spew.Sdump(sugResp))
}

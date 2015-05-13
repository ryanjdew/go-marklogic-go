package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	marklogic "github.com/ryanjdew/go-marklogic-go"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	search "github.com/ryanjdew/go-marklogic-go/search"
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
		authType = marklogic.BasicAuth
	} else if auth == "digest" {
		authType = marklogic.DigestAuth
	} else {
		authType = marklogic.None
	}
	client, err := marklogic.NewClient(host, port, username, password, authType)
	fmt.Print("Client:\n")
	fmt.Print(spew.Sdump(client))
	if err != nil {
		log.Fatal(err)
	}
	query := search.Query{Format: handle.XML}
	query.Queries = []interface{}{
		search.TermQuery{
			Terms: []string{queryStr},
		},
	}

	qh := search.QueryHandle{}
	qh.Decode(query)
	fmt.Print("decoded query:\n")
	fmt.Print(spew.Sdump(qh.Serialized()))
	respHandle := search.ResponseHandle{}
	err = client.StructuredSearch(&qh, 1, 10, &respHandle)
	resp := respHandle.Get()
	fmt.Print("decoded response:\n")
	fmt.Print(spew.Sdump(resp))
	sugRespHandle := search.SuggestionsResponseHandle{}
	err = client.StructuredSuggestions(&qh, queryStr, 10, "", &sugRespHandle)
	sugResp := sugRespHandle.Serialized()
	fmt.Print("decoded response:\n")
	fmt.Print(spew.Sdump(sugResp))
}

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	marklogic "github.com/cchatfield/go-marklogic-go"
	handle "github.com/cchatfield/go-marklogic-go/handle"
	search "github.com/cchatfield/go-marklogic-go/search"
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
	query := search.Query{}
	query.Queries = []interface{}{
		search.TermQuery{
			Terms: []string{queryStr},
		},
	}

	qh := search.QueryHandle{Format: handle.XML}
	qh.Serialize(query)
	fmt.Print("Serialized query:\n")
	fmt.Print(spew.Sdump(qh.Serialized()))
	respHandle := search.ResponseHandle{}
	_ = client.Search().StructuredSearch(&qh, 1, 10, nil, &respHandle)
	resp := respHandle.Get()
	fmt.Print("Serialized response:\n")
	fmt.Print(spew.Sdump(resp))
	sugRespHandle := search.SuggestionsResponseHandle{}
	_ = client.Search().StructuredSuggestions(&qh, queryStr, 10, "", nil, &sugRespHandle)
	sugResp := sugRespHandle.Serialized()
	fmt.Print("Serialized response:\n")
	fmt.Print(spew.Sdump(sugResp))
}

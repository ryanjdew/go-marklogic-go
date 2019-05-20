package alert

import (
	"net/http"

	"github.com/ryanjdew/go-marklogic-go/clients"
	"github.com/ryanjdew/go-marklogic-go/documents"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/search"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

// RulesResponse represents a response from the alert/rules API
type RulesResponse struct {
	Rules []Rule `xml:"http://marklogic.com/rest-api rule" json:"rule,omitempty"`
}

// Rule represents an alert rule from the alert API
type Rule struct {
	Name         string               `xml:"http://marklogic.com/rest-api name" json:"name,omitempty"`
	Description  string               `xml:"http://marklogic.com/rest-api description" json:"description,omitempty"`
	Query        search.CombinedQuery `xml:"http://marklogic.com/appservices/search search" json:"search,omitempty"`
	RuleMetadata RuleMetadata         `xml:"http://marklogic.com/rest-api rule-metadata" json:"rule-metadata,omitempty"`
}

// RuleMetadata represents the metadata in an alert rule from the alert API
type RuleMetadata struct {
	Meta []interface{} `xml:",any" json:",any"`
}

func matchDocument(c *clients.Client, documentDescription documents.DocumentDescription, params map[string]string, response handle.ResponseHandle) error {
	paramsStr := util.MappedParameters("?", "", params)
	paramsStr = util.AddDatabaseParam(paramsStr, c)
	req, err := http.NewRequest("POST", c.Base()+"/alert"+paramsStr, documentDescription.Content)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func matchQuery(c *clients.Client, query handle.Handle, params map[string]string, response handle.ResponseHandle) error {
	paramsStr := util.MappedParameters("?", "", params)
	paramsStr = util.AddDatabaseParam(paramsStr, c)
	req, err := util.BuildRequestFromHandle(c, "POST", c.Base()+"/alert"+paramsStr, query)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func listRules(c *clients.Client, params map[string]string, response handle.ResponseHandle) error {
	paramsStr := util.MappedParameters("?", "", params)
	paramsStr = util.AddDatabaseParam(paramsStr, c)
	req, err := http.NewRequest("GET", c.Base()+"/alert/rules"+paramsStr, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func getRule(c *clients.Client, ruleName string, params map[string]string, response handle.ResponseHandle) error {
	paramsStr := util.MappedParameters("?", "", params)
	paramsStr = util.AddDatabaseParam(paramsStr, c)
	req, err := http.NewRequest("GET", c.Base()+"/alert/rules"+ruleName+paramsStr, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func addRule(c *clients.Client, ruleName string, rule handle.Handle, params map[string]string, response handle.ResponseHandle) error {
	paramsStr := util.MappedParameters("?", "", params)
	paramsStr = util.AddDatabaseParam(paramsStr, c)
	req, err := util.BuildRequestFromHandle(c, "PUT", "/alert/rules/"+ruleName+paramsStr, rule)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func removeRule(c *clients.Client, ruleName string, params map[string]string, response handle.ResponseHandle) error {
	paramsStr := util.MappedParameters("?", "", params)
	paramsStr = util.AddDatabaseParam(paramsStr, c)
	req, err := http.NewRequest("DELETE", c.Base()+"/alert/rules/"+ruleName+paramsStr, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

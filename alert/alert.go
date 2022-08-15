package alert

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/ryanjdew/go-marklogic-go/clients"
	"github.com/ryanjdew/go-marklogic-go/documents"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/search"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

// RulesResponseHandle handles a response from the alert/rules API
type RulesResponseHandle struct {
	*bytes.Buffer
	Format    int
	response  RulesResponse
	timestamp string
}

// GetFormat returns int that represents XML or JSON
func (rh *RulesResponseHandle) GetFormat() int {
	return rh.Format
}

func (rh *RulesResponseHandle) resetBuffer() {
	if rh.Buffer == nil {
		rh.Buffer = new(bytes.Buffer)
	}
	rh.Reset()
}

// Deserialize returns Response struct that represents XML or JSON
func (rh *RulesResponseHandle) Deserialize(bytes []byte) {
	rh.resetBuffer()
	rh.Write(bytes)
	if rh.GetFormat() == handle.JSON {
		json.Unmarshal(bytes, rh)
	} else {
		xml.Unmarshal(bytes, rh)
	}
}

// Deserialized returns deserialised RestartResponse as interface{}
func (rh *RulesResponseHandle) Deserialized() interface{} {
	return &rh.response
}

// AcceptResponse handles an *http.Response
func (rh *RulesResponseHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(rh, resp)
}

// Serialize returns []byte of XML or JSON that represents the Response struct
func (rh *RulesResponseHandle) Serialize(response interface{}) {
	switch response := response.(type) {
	case *RulesResponse:
		rh.response = *(response)
	case RulesResponse:
		rh.response = response
	}
	rh.resetBuffer()
	if rh.GetFormat() == handle.JSON {
		enc := json.NewEncoder(rh.Buffer)
		enc.Encode(&rh.response)
	} else {
		enc := xml.NewEncoder(rh.Buffer)
		enc.Encode(&rh.response)
	}
}

// Get returns deserialised RestartResponse
func (rh *RulesResponseHandle) Get() *RulesResponse {
	return &rh.response
}

// Serialized returns string of XML or JSON
func (rh *RulesResponseHandle) Serialized() string {
	rh.Serialize(rh.response)
	return rh.String()
}

// SetTimestamp sets the timestamp
func (rh *RulesResponseHandle) SetTimestamp(timestamp string) {
	rh.timestamp = timestamp
}

// Timestamp retieves a timestamp
func (rh *RulesResponseHandle) Timestamp() string {
	return rh.timestamp
}

// RulesResponse represents a response from the alert/rules API
type RulesResponse struct {
	XMLName xml.Name `xml:"http://marklogic.com/rest-api rules" json:"rules,omitempty"`
	Rules   []Rule   `xml:"rule" json:"rule,omitempty"`
}

// Rule represents an alert rule from the alert API
type Rule struct {
	Name         string                     `xml:"name" json:"name,omitempty"`
	Description  string                     `xml:"description" json:"description,omitempty"`
	Query        search.CombinedQuery       `xml:"http://marklogic.com/appservices/search search" json:"search,omitempty"`
	RuleMetadata util.SerializableStringMap `xml:"rule-metadata" json:"rule-metadata,omitempty"`
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

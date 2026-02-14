package eval

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Service provides an interface for MarkLogic /v1/eval endpoint operations.
// The /v1/eval endpoint allows ad-hoc execution of server-side code (XQuery or JavaScript).
type Service struct {
	client *clients.Client
}

// NewService returns a Service for managing eval operations.
func NewService(c *clients.Client) *Service {
	return &Service{client: c}
}

// EvalXQuery evaluates XQuery code on the server.
// Parameters:
//   - code: The XQuery code to execute
//   - params: Optional map of external variable definitions (e.g., "var" -> "{ \"type\": \"xs:string\", \"value\": \"test\" }")
//   - response: ResponseHandle to deserialize the response
func (s *Service) EvalXQuery(code string, params map[string]string, response handle.ResponseHandle) error {
	return evalCode(s.client, "xquery", code, params, response)
}

// EvalJavaScript evaluates JavaScript code on the server.
// Parameters:
//   - code: The JavaScript code to execute
//   - params: Optional map of external variable definitions
//   - response: ResponseHandle to deserialize the response
func (s *Service) EvalJavaScript(code string, params map[string]string, response handle.ResponseHandle) error {
	return evalCode(s.client, "javascript", code, params, response)
}

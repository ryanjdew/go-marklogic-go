package management

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	handle "github.com/cchatfield/go-marklogic-go/handle"
	test "github.com/cchatfield/go-marklogic-go/test"
	textHelper "github.com/cchatfield/go-marklogic-go/test/text"
)

var serverPropertiesWantResp = `{
  "server-name": "samplestack",
  "group-name": "Default",
  "server-type": "http",
  "enabled": true,
  "root": "/",
  "port": 8006,
  "webDAV": false,
  "execute": true,
  "display-last-login": false,
  "address": "0.0.0.0",
  "backlog": 512,
  "threads": 32,
  "request-timeout": 30,
  "keep-alive-timeout": 5,
  "session-timeout": 3600,
  "max-time-limit": 3600,
  "default-time-limit": 600,
  "max-inference-size": 500,
  "default-inference-size": 100,
  "static-expires": 3600,
  "pre-commit-trigger-depth": 1000,
  "pre-commit-trigger-limit": 10000,
  "collation": "http://marklogic.com/collation/",
  "authentication": "digest",
  "internal-security": true,
  "concurrent-request-limit": 0,
  "compute-content-length": true,
  "log-errors": false,
  "debug-allow": true,
  "profile-allow": true,
  "default-xquery-version": "1.0-ml",
  "multi-version-concurrency-control": "contemporaneous",
  "distribute-timestamps": "fast",
  "output-sgml-character-entities": "none",
  "output-encoding": "UTF-8",
  "output-method": "default",
  "output-byte-order-mark": "default",
  "output-cdata-section-namespace-uri": "",
  "output-cdata-section-localname": "",
  "output-doctype-public": "",
  "output-doctype-system": "",
  "output-escape-uri-attributes": "default",
  "output-include-content-type": "default",
  "output-indent": "default",
  "output-indent-untyped": "default",
  "output-media-type": "",
  "output-normalization-form": "none",
  "output-omit-xml-declaration": "default",
  "output-standalone": "omit",
  "output-undeclare-prefixes": "default",
  "output-version": "",
  "output-include-default-attributes": "default",
  "default-error-format": "json",
  "error-handler": "/MarkLogic/rest-api/error-handler.xqy",
  "url-rewriter": "/MarkLogic/rest-api/rewriter.xml",
  "rewrite-resolves-globally": true,
  "ssl-certificate-template": 0,
  "ssl-allow-sslv3": true,
  "ssl-allow-tls": true,
  "ssl-hostname": "",
  "ssl-ciphers": "ALL:!LOW:@STRENGTH",
  "ssl-require-client-certificate": true,
  "content-database": "samplestack",
  "modules-database": "samplestack-modules",
  "default-user": "nobody"
}`

func TestServerProperties(t *testing.T) {
	client, server := test.ManagementClient(serverPropertiesWantResp)
	defer server.Close()
	propertiesHandle := ServerPropertiesHandle{Format: handle.JSON}
	GetServerProperties(client, "Documents", "Default", &propertiesHandle)
	want := textHelper.NormalizeSpace(serverPropertiesWantResp)
	result := textHelper.NormalizeSpace(propertiesHandle.Serialized())
	if !reflect.DeepEqual(want, result) {
		t.Errorf("Server Properties Results = %+v, Want = %+v", spew.Sdump(result), spew.Sdump(want))
	}
}

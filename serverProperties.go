package goMarklogicGo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// ServerProperties represents the properties of a MarkLogic AppServer
type ServerProperties struct {
	ServerName                     string `xml:"http://marklogic.com/manage server-name" json:"server-name"`
	GroupName                      string `xml:"http://marklogic.com/manage group-name" json:"group-name"`
	ServerType                     string `xml:"http://marklogic.com/manage server-type" json:"server-type"`
	Enabled                        bool   `xml:"http://marklogic.com/manage enabled" json:"enabled"`
	Root                           string `xml:"http://marklogic.com/manage root" json:"root"`
	Port                           int    `xml:"http://marklogic.com/manage port" json:"port"`
	WebDAV                         bool   `xml:"http://marklogic.com/manage webDAV" json:"webDAV"`
	Execute                        bool   `xml:"http://marklogic.com/manage execute" json:"execute"`
	DisplayLastLogin               bool   `xml:"http://marklogic.com/manage display-last-login" json:"display-last-login"`
	Address                        string `xml:"http://marklogic.com/manage address" json:"address"`
	Backlog                        int    `xml:"http://marklogic.com/manage backlog" json:"backlog"`
	Threads                        int    `xml:"http://marklogic.com/manage threads" json:"threads"`
	RequestTimeout                 int    `xml:"http://marklogic.com/manage request-timeout" json:"request-timeout"`
	KeepAliveTimeout               int    `xml:"http://marklogic.com/manage keep-alive-timeout" json:"keep-alive-timeout"`
	SessionTimeout                 int    `xml:"http://marklogic.com/manage session-timeout" json:"session-timeout"`
	MaxTimeLimit                   int    `xml:"http://marklogic.com/manage max-time-limit" json:"max-time-limit"`
	DefaultTimeLimit               int    `xml:"http://marklogic.com/manage default-time-limit" json:"default-time-limit"`
	MaxInferenceSize               int    `xml:"http://marklogic.com/manage max-inference-size" json:"max-inference-size"`
	DefaultInferenceSize           int    `xml:"http://marklogic.com/manage default-inference-size" json:"default-inference-size"`
	StaticExpires                  int    `xml:"http://marklogic.com/manage static-expires" json:"static-expires"`
	PreCommitTriggerDepth          int    `xml:"http://marklogic.com/manage pre-commit-trigger-depth" json:"pre-commit-trigger-depth"`
	PreCommitTriggerLimit          int    `xml:"http://marklogic.com/manage pre-commit-trigger-limit" json:"pre-commit-trigger-limit"`
	Collation                      string `xml:"http://marklogic.com/manage collation" json:"collation"`
	Authentication                 string `xml:"http://marklogic.com/manage authentication" json:"authentication"`
	InternalSecurity               bool   `xml:"http://marklogic.com/manage internal-security" json:"internal-security"`
	ConcurrentRequestLimit         int    `xml:"http://marklogic.com/manage concurrent-request-limit" json:"concurrent-request-limit"`
	ComputeContentLength           bool   `xml:"http://marklogic.com/manage compute-content-length" json:"compute-content-length"`
	LogErrors                      bool   `xml:"http://marklogic.com/manage log-errors" json:"log-errors"`
	DebugAllow                     bool   `xml:"http://marklogic.com/manage debug-allow" json:"debug-allow"`
	ProfileAllow                   bool   `xml:"http://marklogic.com/manage profile-allow" json:"profile-allow"`
	DefaultXqueryVersion           string `xml:"http://marklogic.com/manage default-xquery-version" json:"default-xquery-version"`
	MultiVersionConcurrencyControl string `xml:"http://marklogic.com/manage multi-version-concurrency-control" json:"multi-version-concurrency-control"`
	DistributeTimestamps           string `xml:"http://marklogic.com/manage distribute-timestamps" json:"distribute-timestamps"`
	OutputSgmlCharacterEntities    string `xml:"http://marklogic.com/manage output-sgml-character-entities" json:"output-sgml-character-entities"`
	OutputEncoding                 string `xml:"http://marklogic.com/manage output-encoding" json:"output-encoding"`
	OutputMethod                   string `xml:"http://marklogic.com/manage output-method" json:"output-method"`
	OutputByteOrderMark            string `xml:"http://marklogic.com/manage output-byte-order-mark" json:"output-byte-order-mark"`
	OutputCdataSectionNamespaceURI string `xml:"http://marklogic.com/manage output-cdata-section-namespace-uri" json:"output-cdata-section-namespace-uri"`
	OutputCdataSectionLocalname    string `xml:"http://marklogic.com/manage output-cdata-section-localname" json:"output-cdata-section-localname"`
	OutputDoctypePublic            string `xml:"http://marklogic.com/manage output-doctype-public" json:"output-doctype-public"`
	OutputDoctypeSystem            string `xml:"http://marklogic.com/manage output-doctype-system" json:"output-doctype-system"`
	OutputEscapeURIAttributes      string `xml:"http://marklogic.com/manage output-escape-uri-attributes" json:"output-escape-uri-attributes"`
	OutputIncludeContentType       string `xml:"http://marklogic.com/manage output-include-content-type" json:"output-include-content-type"`
	OutputIndent                   string `xml:"http://marklogic.com/manage output-indent" json:"output-indent"`
	OutputIndentUntyped            string `xml:"http://marklogic.com/manage output-indent-untyped" json:"output-indent-untyped"`
	OutputMediaType                string `xml:"http://marklogic.com/manage output-media-type" json:"output-media-type"`
	OutputNormalizationForm        string `xml:"http://marklogic.com/manage output-normalization-form" json:"output-normalization-form"`
	OutputOmitXMLDeclaration       string `xml:"http://marklogic.com/manage output-omit-xml-declaration" json:"output-omit-xml-declaration"`
	OutputStandalone               string `xml:"http://marklogic.com/manage output-standalone" json:"output-standalone"`
	OutputUndeclarePrefixes        string `xml:"http://marklogic.com/manage output-undeclare-prefixes" json:"output-undeclare-prefixes"`
	OutputVersion                  string `xml:"http://marklogic.com/manage output-version" json:"output-version"`
	OutputIncludeDefaultAttributes string `xml:"http://marklogic.com/manage output-include-default-attributes" json:"output-include-default-attributes"`
	DefaultErrorFormat             string `xml:"http://marklogic.com/manage default-error-format" json:"default-error-format"`
	ErrorHandler                   string `xml:"http://marklogic.com/manage error-handler" json:"error-handler"`
	URLRewriter                    string `xml:"http://marklogic.com/manage url-rewriter" json:"url-rewriter"`
	RewriteResolvesGlobally        bool   `xml:"http://marklogic.com/manage rewrite-resolves-globally" json:"rewrite-resolves-globally"`
	SslCertificateTemplate         int    `xml:"http://marklogic.com/manage ssl-certificate-template" json:"ssl-certificate-template"`
	SslAllowSslv3                  bool   `xml:"http://marklogic.com/manage ssl-allow-sslv3" json:"ssl-allow-sslv3"`
	SslAllowTLS                    bool   `xml:"http://marklogic.com/manage ssl-allow-tls" json:"ssl-allow-tls"`
	SslHostname                    string `xml:"http://marklogic.com/manage ssl-hostname" json:"ssl-hostname"`
	SslCiphers                     string `xml:"http://marklogic.com/manage ssl-ciphers" json:"ssl-ciphers"`
	SslRequireClientCertificate    bool   `xml:"http://marklogic.com/manage ssl-require-client-certificate" json:"ssl-require-client-certificate"`
	ContentDatabase                string `xml:"http://marklogic.com/manage content-database" json:"content-database"`
	ModulesDatabase                string `xml:"http://marklogic.com/manage modules-database" json:"modules-database"`
	DefaultUser                    string `xml:"http://marklogic.com/manage default-user" json:"default-user"`
}

// SetServerProperties sets the database properties
func (mc *ManagementClient) SetServerProperties(serverName string, propertiesHandle Handle) error {
	var reqType string
	if propertiesHandle.GetFormat() == JSON {
		reqType = "json"
	} else {
		reqType = "xml"
	}
	buf := new(bytes.Buffer)
	buf.Write([]byte(propertiesHandle.Serialized()))
	req, _ := http.NewRequest("PUT", mc.Base()+"/servers/"+serverName+"/properties?format="+reqType, buf)
	applyAuth(mc, req)
	req.Header.Add("Content-Type", "application/"+reqType)
	resp, err := mc.HTTPClient().Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

// GetServerProperties sets the database properties
func (mc *ManagementClient) GetServerProperties(serverName string, propertiesHandle Handle) error {
	var reqType string
	if propertiesHandle.GetFormat() == JSON {
		reqType = "json"
	} else {
		reqType = "xml"
	}
	req, _ := http.NewRequest("GET", mc.Base()+"/servers/"+serverName+"/properties?format="+reqType, nil)
	applyAuth(mc, req)
	req.Header.Add("Content-Type", "application/"+reqType)
	resp, err := mc.HTTPClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	propertiesHandle.Encode([]byte(contents))
	return err
}

// ServerPropertiesHandle is a handle that places the results into
// a ServerProperties struct
type ServerPropertiesHandle struct {
	Format           int
	bytes            []byte
	serverProperties *ServerProperties
}

// GetFormat returns int that represents XML or JSON
func (sh *ServerPropertiesHandle) GetFormat() int {
	return sh.Format
}

// Encode returns Query struct that represents XML or JSON
func (sh *ServerPropertiesHandle) Encode(bytes []byte) *ServerProperties {
	sh.bytes = bytes
	sh.serverProperties = &ServerProperties{}
	if sh.GetFormat() == JSON {
		json.Unmarshal(bytes, &sh.serverProperties)
	} else {
		xml.Unmarshal(bytes, &sh.serverProperties)
	}
	return sh.serverProperties
}

// Decode returns []byte of XML or JSON that represents the Query struct
func (sh *ServerPropertiesHandle) Decode(serverProperties *ServerProperties) []byte {
	sh.serverProperties = serverProperties
	buf := new(bytes.Buffer)
	if sh.GetFormat() == JSON {
		enc := json.NewEncoder(buf)
		enc.Encode(sh.serverProperties)
	} else {
		enc := xml.NewEncoder(buf)
		enc.Encode(sh.serverProperties)
	}
	sh.bytes = buf.Bytes()
	return sh.bytes
}

// Get returns string of XML or JSON
func (sh *ServerPropertiesHandle) Get() *ServerProperties {
	return sh.serverProperties
}

// Serialized returns string of XML or JSON
func (sh *ServerPropertiesHandle) Serialized() string {
	buf := new(bytes.Buffer)
	if sh.GetFormat() == JSON {
		enc := json.NewEncoder(buf)
		enc.Encode(sh.serverProperties)
	} else {
		enc := xml.NewEncoder(buf)
		enc.Encode(sh.serverProperties)
	}
	sh.bytes = buf.Bytes()
	return string(sh.bytes)
}

package documents

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

// DocumentDescription describes a document to write
type DocumentDescription struct {
	URI       string
	Content   io.ReadWriter
	Metadata  *Metadata
	Format    int
	VersionID int
}

// GetFormat returns int that represents XML or JSON
func (dd *DocumentDescription) GetFormat() int {
	return dd.Format
}

func toURIs(docs []*DocumentDescription) []string {
	uris := []string{}
	for _, doc := range docs {
		uris = append(uris, doc.URI)
	}
	return uris
}

func read(c *clients.Client, uris []string, categories []string, transform *util.Transform, transaction *util.Transaction, response handle.ResponseHandle) error {
	params := buildParameters(uris, categories, nil, nil, nil, transform)
	params = util.AddDatabaseParam(params, c)
	params = util.AddTransactionParam(params, transaction)
	req, err := http.NewRequest("GET", c.Base()+"/documents"+params, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func write(c *clients.Client, documents []*DocumentDescription, transform *util.Transform, transaction *util.Transaction, response handle.ResponseHandle) error {
	channel := make(chan error)
	var errReturn error
	for _, doc := range documents {
		go func(doc *DocumentDescription) {
			metadata := doc.Metadata
			params := buildParameters([]string{doc.URI}, nil, metadata.Collections, metadata.PermissionsMap(), metadata.Properties, transform)
			params = util.AddDatabaseParam(params, c)
			params = util.AddTransactionParam(params, transaction)
			req, err := http.NewRequest("PUT", c.Base()+"/documents"+params, doc.Content)
			if err == nil {
				err = util.Execute(c, req, response)
			}
			channel <- err
		}(doc)
	}
	for range documents {
		if errReturn == nil {
			errReturn = <-channel
		} else {
			<-channel
		}
	}
	return errReturn
}

func writeSet(c *clients.Client, documents []*DocumentDescription, metadata handle.Handle, transform *util.Transform, transaction *util.Transaction, response handle.ResponseHandle) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	params := ""
	if transform != nil {
		params = "?" + transform.ToParameters()
	}
	params = util.AddDatabaseParam(params, c)
	params = util.AddTransactionParam(params, transaction)
	metaHeader := &textproto.MIMEHeader{}
	metadataSerialized := metadata.Serialized()
	if metadata.GetFormat() == handle.JSON {
		metaHeader.Add("Content-Type", "application/json")
	} else if metadata.GetFormat() == handle.JSON {
		metaHeader.Add("Content-Type", "application/xml")
	} else if metadata.GetFormat() == handle.TEXTPLAIN {
		metaHeader.Add("Content-Type", "text/plain")
	} else {
		metaHeader.Add("Content-Type", "application/octet-stream")
	}
	metaHeader.Add("Content-Disposition", "inline; category=metadata")
	metaHeader.Add("Content-Length", strconv.Itoa(len(metadataSerialized)))
	writer.CreatePart(*metaHeader)
	body.Write([]byte(metadataSerialized))
	for _, doc := range documents {
		if doc.Metadata != nil {
			docMetaHeader := &textproto.MIMEHeader{}
			metadataHandle := &MetadataHandle{metadata: *doc.Metadata}
			docMetadataSerialized := metadataHandle.Serialized()
			if metadataHandle.GetFormat() == handle.XML {
				docMetaHeader.Add("Content-Type", "application/xml")
			} else {
				docMetaHeader.Add("Content-Type", "application/json")
			}
			docMetaHeader.Add("Content-Disposition", "inline; category=metadata; filename=\""+doc.URI+"\"")
			docMetaHeader.Add("Content-Length", strconv.Itoa(len(docMetadataSerialized)))
			writer.CreatePart(*docMetaHeader)
			body.Write([]byte(docMetadataSerialized))
		}
		var docContentBytes []byte
		if doc.Content != nil {
			docContentBytes, _ = ioutil.ReadAll(doc.Content)
		}
		doc.Content = bytes.NewBuffer(docContentBytes)
		header := &textproto.MIMEHeader{}
		if doc.GetFormat() == handle.JSON {
			header.Add("Content-Type", "application/json")
		} else if doc.GetFormat() == handle.JSON {
			header.Add("Content-Type", "application/xml")
		} else if doc.GetFormat() == handle.TEXTPLAIN {
			header.Add("Content-Type", "text/plain")
		} else {
			header.Add("Content-Type", "application/octet-stream")
		}
		header.Add("Content-Disposition", "attachment; filename=\""+doc.URI+"\"")
		header.Add("Content-Length", strconv.Itoa(len(docContentBytes)))
		writer.CreatePart(*header)
		body.Write(docContentBytes)
	}
	body.Write([]byte("\r\n--" + writer.Boundary() + "--"))
	req, err := http.NewRequest("POST", c.Base()+"/documents"+params, body)
	req.Header.Add("Content-Type", "multipart/mixed; boundary="+writer.Boundary())
	if err == nil {
		err = util.Execute(c, req, response)
	}
	return err
}

func buildParameters(uris []string, categories []string, collections []string, permissions map[string]string, properties map[string]string, transform *util.Transform) string {
	params := "?"
	params = util.RepeatingParameters(params, "uri", uris)
	params = util.RepeatingParameters(params, "category", categories)
	params = util.RepeatingParameters(params, "collection", collections)
	params = util.MappedParameters(params, "perm", permissions)
	params = util.MappedParameters(params, "prop", properties)
	if transform != nil {
		params = params + transform.ToParameters()
	}
	return params
}

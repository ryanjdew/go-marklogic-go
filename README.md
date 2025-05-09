Go, MarkLogic. Go!
=========

This a Go library for interacting with MarkLogic's REST APIs.

Go, MarkLogic. Go!  
Look at those MarkLogicians go!  
Where is MarkLogic going with those programming stars?  
What are they going to do?  
Where are those MarkLogicians going?  
Look where they are going.  
They are all going to that big data out there.  
It's a data party. A BIG data party!  
XML data. JSON data. Semantic data.  

Status
=========
[![GoDoc](https://godoc.org/github.com/ryanjdew/go-marklogic-go?status.svg)](https://godoc.org/github.com/ryanjdew/go-marklogic-go) [![Build Status](https://cloud.drone.io/api/badges/ryanjdew/go-marklogic-go/status.svg)](https://cloud.drone.io/ryanjdew/go-marklogic-go)

Sample Code
=========

```go
import (
	"fmt"
	marklogic "github.com/ryanjdew/go-marklogic-go"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	search "github.com/ryanjdew/go-marklogic-go/search"
)
func main() {
  db, _ := marklogic.NewClient("localhost", 8050, "admin", "admin", marklogic.DigestAuth)
  query := search.Query{}
  query.Queries = []any{
    search.TermQuery{
      Terms: []string{queryStr},
    },
  }
  qh := search.QueryHandle{Format: handle.XML}
  qh.Serialize(query)
  respHandle := search.ResponseHandle{}
  err = db.Search().StructuredSearch(&qh, 1, 10, nil, &respHandle)
  resp := respHandle.Get()
  fmt.Print(respHandle.Serialized())
}
```

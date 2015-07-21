package config

import (
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/test"
	"github.com/ryanjdew/go-marklogic-go/test/text"
)

func TestListExtensions(t *testing.T) {
	want := testHelper.NormalizeSpace(`{
  "resources": {
    "resource": [
      {
        "name": "example",
        "title": "",
        "version": "",
        "provider-name": "Acme Widgets",
        "description": "",
        "methods": {
          "method": [
            {
              "parameter": [
                {
                  "parameter-name": "the-uri",
                  "parameter-type": "string"
                }
              ],
              "method-name": "get"
            }
          ]
        },
        "resource-source": "/v1/resources/example"
      }
    ]
  }
}`)
	client, _ := test.Client(want)
	responseHandle := handle.MapHandle{
		Format: handle.JSON,
	}
	err := listExtensions(client, "/", &responseHandle)
	result := testHelper.NormalizeSpace(responseHandle.Serialized())
	if err != nil {
		t.Errorf("Error encountered: %v", err)
	} else if want != result {
		t.Errorf("Result = %v, want = %v", result, want)
	}
}

func TestDeleteExtensions(t *testing.T) {
	want := testHelper.NormalizeSpace(``)
	client, _ := test.Client(want)
	err := deleteExtensions(client, "/")
	if err != nil {
		t.Errorf("Error encountered: %v", err)
	}
}

func TestCreateExtension(t *testing.T) {
	want := testHelper.NormalizeSpace(`{
  "resources": {
    "resource": [
      {
        "name": "example",
        "title": "",
        "version": "",
        "provider-name": "Acme Widgets",
        "description": "",
        "methods": {
          "method": [
            {
              "parameter": [
                {
                  "parameter-name": "the-uri",
                  "parameter-type": "string"
                }
              ],
              "method-name": "get"
            }
          ]
        },
        "resource-source": "/v1/resources/example"
      }
    ]
  }
}`)
	client, _ := test.Client(want)
	responseHandle := handle.MapHandle{
		Format: handle.JSON,
	}
	err := listExtensions(client, "/", &responseHandle)
	result := testHelper.NormalizeSpace(responseHandle.Serialized())
	if err != nil {
		t.Errorf("Error encountered: %v", err)
	} else if want != result {
		t.Errorf("Result = %v, want = %v", result, want)
	}
}

func TestListResources(t *testing.T) {
	want := testHelper.NormalizeSpace(`{
  "resources": {
    "resource": [
      {
        "name": "example",
        "title": "",
        "version": "",
        "provider-name": "Acme Widgets",
        "description": "",
        "methods": {
          "method": [
            {
              "parameter": [
                {
                  "parameter-name": "the-uri",
                  "parameter-type": "string"
                }
              ],
              "method-name": "get"
            }
          ]
        },
        "resource-source": "/v1/resources/example"
      }
    ]
  }
}`)
	client, _ := test.Client(want)
	responseHandle := handle.MapHandle{
		Format: handle.JSON,
	}
	err := listExtensions(client, "/", &responseHandle)
	result := testHelper.NormalizeSpace(responseHandle.Serialized())
	if err != nil {
		t.Errorf("Error encountered: %v", err)
	} else if want != result {
		t.Errorf("Result = %v, want = %v", result, want)
	}
}

func TestGetResourceInfo(t *testing.T) {
	want := testHelper.NormalizeSpace(`{
  "resources": {
    "resource": [
      {
        "name": "example",
        "title": "",
        "version": "",
        "provider-name": "Acme Widgets",
        "description": "",
        "methods": {
          "method": [
            {
              "parameter": [
                {
                  "parameter-name": "the-uri",
                  "parameter-type": "string"
                }
              ],
              "method-name": "get"
            }
          ]
        },
        "resource-source": "/v1/resources/example"
      }
    ]
  }
}`)
	client, _ := test.Client(want)
	responseHandle := handle.MapHandle{
		Format: handle.JSON,
	}
	err := listExtensions(client, "/", &responseHandle)
	result := testHelper.NormalizeSpace(responseHandle.Serialized())
	if err != nil {
		t.Errorf("Error encountered: %v", err)
	} else if want != result {
		t.Errorf("Result = %v, want = %v", result, want)
	}
}

func TestCreateResource(t *testing.T) {
	want := testHelper.NormalizeSpace(`{
  "resources": {
    "resource": [
      {
        "name": "example",
        "title": "",
        "version": "",
        "provider-name": "Acme Widgets",
        "description": "",
        "methods": {
          "method": [
            {
              "parameter": [
                {
                  "parameter-name": "the-uri",
                  "parameter-type": "string"
                }
              ],
              "method-name": "get"
            }
          ]
        },
        "resource-source": "/v1/resources/example"
      }
    ]
  }
}`)
	client, _ := test.Client(want)
	responseHandle := handle.MapHandle{
		Format: handle.JSON,
	}
	err := listExtensions(client, "/", &responseHandle)
	result := testHelper.NormalizeSpace(responseHandle.Serialized())
	if err != nil {
		t.Errorf("Error encountered: %v", err)
	} else if want != result {
		t.Errorf("Result = %v, want = %v", result, want)
	}
}

func TestDeleteResource(t *testing.T) {
	want := testHelper.NormalizeSpace(`{
  "resources": {
    "resource": [
      {
        "name": "example",
        "title": "",
        "version": "",
        "provider-name": "Acme Widgets",
        "description": "",
        "methods": {
          "method": [
            {
              "parameter": [
                {
                  "parameter-name": "the-uri",
                  "parameter-type": "string"
                }
              ],
              "method-name": "get"
            }
          ]
        },
        "resource-source": "/v1/resources/example"
      }
    ]
  }
}`)
	client, _ := test.Client(want)
	responseHandle := handle.MapHandle{
		Format: handle.JSON,
	}
	err := listExtensions(client, "/", &responseHandle)
	result := testHelper.NormalizeSpace(responseHandle.Serialized())
	if err != nil {
		t.Errorf("Error encountered: %v", err)
	} else if want != result {
		t.Errorf("Result = %v, want = %v", result, want)
	}
}

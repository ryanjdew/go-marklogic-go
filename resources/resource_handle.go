package resources

import (
	"bytes"
	"encoding/json"
	"net/http"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// ResourceExtension represents a user-defined REST resource extension
type ResourceExtension struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	Provider string `json:"provider"`
	URI      string `json:"uri"`
}

// ResourceExtensionHandle serializes/deserializes resource extension metadata
type ResourceExtensionHandle struct {
	*bytes.Buffer
	Format    int
	resource  ResourceExtension
	timestamp string
}

func (h *ResourceExtensionHandle) GetFormat() int {
	return h.Format
}

func (h *ResourceExtensionHandle) resetBuffer() {
	if h.Buffer == nil {
		h.Buffer = new(bytes.Buffer)
	}
	h.Reset()
}

func (h *ResourceExtensionHandle) Serialize(v interface{}) {
	h.resource = v.(ResourceExtension)
	h.resetBuffer()
	if h.GetFormat() == handle.JSON {
		enc := json.NewEncoder(h.Buffer)
		enc.Encode(h.resource)
	}
}

func (h *ResourceExtensionHandle) Deserialize(b []byte) {
	h.resetBuffer()
	h.Write(b)
	h.resource = ResourceExtension{}
	if h.GetFormat() == handle.JSON {
		json.Unmarshal(b, &h.resource)
	}
}

func (h *ResourceExtensionHandle) Deserialized() interface{} {
	return h.resource
}

func (h *ResourceExtensionHandle) Serialized() string {
	h.Serialize(h.resource)
	return h.String()
}

func (h *ResourceExtensionHandle) SetTimestamp(ts string) {
	h.timestamp = ts
}

func (h *ResourceExtensionHandle) Timestamp() string {
	return h.timestamp
}

func (h *ResourceExtensionHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(h, resp)
}

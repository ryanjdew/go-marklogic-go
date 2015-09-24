package admin

import (
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

func TestTimestampResponseHandleDeserialize(t *testing.T) {
	want := `2013-05-15T10:34:38.932514-07:00`
	result := TimestampResponseHandle{Format: handle.TEXTPLAIN}
	result.Deserialize([]byte(want))
	if result.timestamp != want {
		t.Errorf("Not equal - TimestampResponseHandle timestamp = %+v, Want = %+v", result.timestamp, want)
	}
}

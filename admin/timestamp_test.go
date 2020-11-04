package admin

import (
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	testHelper "github.com/ryanjdew/go-marklogic-go/test"
)

func TestTimestampResponseHandleDeserialize(t *testing.T) {
	want := "2013-05-15T10:34:38.932514-07:00"
	result := &TimestampResponseHandle{Format: handle.TEXTPLAIN}
	testHelper.RoundTripSerialization(t, "TimestampResponseHandle", want, result, want)
}

package server

import (
	"testing"
)

func TestServer(t *testing.T) {
	ListenAndServe("0.0.0.0:9900", "topSecret")
}

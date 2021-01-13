package pics

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestRespond(t *testing.T) {
	type test struct {
		payload string
		code    int
	}

	tests := []test{
		{payload: "foobar", code: 200},
		{payload: "foobar", code: 404},
		{payload: "foobar", code: 400},
	}
	s := NewServer(newDefaultConfig(), nil)

	for _, tc := range tests {
		rr := httptest.NewRecorder()
		s.respond(rr, tc.payload, tc.code)
		if rr.Code != tc.code {
			t.Errorf("expected http status code to be %d instead got %d", tc.code, rr.Code)
		}
		expected := fmt.Sprintf(`"%s"`, tc.payload)
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	}
}

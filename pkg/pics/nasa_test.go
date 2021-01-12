package pics

import (
	"testing"
	"time"
)

func TestGetDays(t *testing.T) {
	span := 24 * 7
	expected := span/24 + 1
	d, err := getJobs(time.Now(), time.Now().Add(time.Hour*24*7))
	if err != nil {
		t.Error(err.Error())
	}

	if len(d) != expected {
		t.Errorf("expected to get %d days instead got %d", expected, len(d))
	}

}

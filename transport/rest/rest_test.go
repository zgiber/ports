package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zgiber/ports/service"
)

func TestUpdateAndListPorts(t *testing.T) {
	testService := &testService{}
	s := NewServer(context.Background(), testService, "")

	// given a valid update request
	update := map[string]*service.PortDetails{
		"BCD": {
			Name:        "B",
			City:        "B",
			Country:     "B",
			Coordinates: []float64{2, 3},
			Province:    "B",
			Timezone:    "B",
			UNLocs:      []string{"BCD"},
			Code:        "2",
		},
		"ABC": {
			Name:        "A",
			City:        "A",
			Country:     "A",
			Coordinates: []float64{1, 2},
			Province:    "A",
			Timezone:    "A",
			UNLocs:      []string{"ABC"},
			Code:        "1",
		},
	}
	b, _ := json.Marshal(update)
	request, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(b))

	// when the request is handled
	rec := httptest.NewRecorder()
	s.postUpdateBatch(rec, request)

	// then the service is called with the correct parameters
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Len(t, testService.update, len(update))

	// (reusing previous test's state for simplicity)

	// given a subsequent list request
	request, _ = http.NewRequest(http.MethodGet, "http://host/update?from_id=A&limit=2", nil)

	// when the request is handled
	rec = httptest.NewRecorder()
	s.getPortList(rec, request)

	// then the service is called with the correct parameters
	assert.EqualValues(t, service.PortsFilter{
		AfterID:  "A",
		MaxItems: 2,
	}, testService.list)

	// and the correct response is returned
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	results := []*service.Port{}
	assert.NoError(t, json.NewDecoder(rec.Body).Decode(&results))
	assert.Len(t, results, len(update))

	for _, res := range results {
		assert.EqualValues(t, &res.Details, update[res.ID])
	}

}

type testService struct {
	update []*service.Port
	list   service.PortsFilter
}

func (ts *testService) UpdatePorts(ctx context.Context, feed service.Feed) error {
	ts.update = []*service.Port{}
	for {
		p, _ := feed.Next()
		if p != nil {
			ts.update = append(ts.update, p)
			continue
		}
		return nil
	}
}

func (ts *testService) ListPorts(ctx context.Context, filter service.PortsFilter) ([]*service.Port, error) {
	ts.list = filter
	return ts.update, nil
}

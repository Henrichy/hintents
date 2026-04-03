// Copyright 2026 Erst Users
// SPDX-License-Identifier: Apache-2.0

package simulator

import (
	"encoding/json"
	"testing"
)

func TestSimulationResponseUnmarshalSnapshots(t *testing.T) {
	payload := []byte(`{
		"status":"ok",
		"snapshots":{
			"inline":{
				"post-run":{
					"ledger_entries":[["a2V5","dmFsdWU="]],
					"linear_memory":"AQID"
				}
			},
			"ids":["snap-2"]
		}
	}`)

	var resp SimulationResponse
	if err := json.Unmarshal(payload, &resp); err != nil {
		t.Fatalf("expected no error unmarshalling snapshots, got %v", err)
	}
	if resp.Snapshots == nil {
		t.Fatal("expected snapshots payload")
	}
	if _, ok := resp.Snapshots.Inline["post-run"]; !ok {
		t.Fatal("expected inline snapshot with key post-run")
	}
	if got := len(resp.Snapshots.IDs); got != 1 {
		t.Fatalf("expected 1 lazy snapshot id, got %d", got)
	}
}

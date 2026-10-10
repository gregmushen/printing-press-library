package cli

import "encoding/json"

// applyTelnyxSyncDefaults fills per-resource query params the Telnyx API
// requires but the spec does not mark as defaults. GET /detail_records
// returns 400 "Record type is missing" without filter[record_type], and the
// messaging records are what conversation, spend, failures and inbox read.
// An explicit --resource-param detail-records:filter[record_type]=... wins.
func applyTelnyxSyncDefaults(p *syncUserParams) {
	if p.perResource == nil {
		p.perResource = map[string]map[string]string{}
	}
	dr := p.perResource["detail-records"]
	if dr == nil {
		dr = map[string]string{}
		p.perResource["detail-records"] = dr
	}
	if _, ok := dr["filter[record_type]"]; !ok {
		dr["filter[record_type]"] = "messaging"
	}
}

// withDefaultRecordType adds record_type=messaging to a --filter JSON object
// when the caller did not choose one. Non-object JSON is returned unchanged
// so the generated parser reports the error.
func withDefaultRecordType(raw string) string {
	if raw == "" {
		return `{"record_type":"messaging"}`
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil || m == nil {
		// Not a JSON object (including literal null): leave it for the
		// generated parser to accept or reject; never index a nil map.
		if m == nil && err == nil {
			return `{"record_type":"messaging"}`
		}
		return raw
	}
	if _, ok := m["record_type"]; ok {
		return raw
	}
	m["record_type"] = "messaging"
	b, err := json.Marshal(m)
	if err != nil {
		return raw
	}
	return string(b)
}

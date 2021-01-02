package defectdojo_client

import (
	"encoding/json"
	"fmt"
)

type Engagement struct {
	EngagementId   int    `json:"id,omitempty"`
	ProductId      int    `json:"product"`
	StartDate      string `json:"target_start"`
	EndDate        string `json:"target_end"`
	EngagementType string `json:"engagement_type"`
	EngagementName string `json:"name"`
}

func (c *DefectdojoClient) CreateEngagement(p *Product, report_type string) (*Engagement, error) {
	engagement_req := Engagement{
		ProductId:      p.Id,
		StartDate:      "2021-01-01",
		EndDate:        "2021-01-01",
		EngagementType: "CI/CD",
		EngagementName: report_type,
	}

	payload, err := c.BuildJsonRequestBytez(engagement_req)
	if err != nil {
		return nil, err
	}
	resp, err := c.DoPost("engagements", payload, APPLICATION_JSON)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e *Engagement
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&e); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err)
	}

	return e, nil
}

func (c *DefectdojoClient) UploadReport(engagement_id int, report_type string, report_bytez []byte) (*Engagement, error) {
	form := map[string]string{
		"engagement": fmt.Sprint(engagement_id),
		"scan_type":  report_type,
	}

	bytez, header, err := c.BuildMultipartFormBytez(form, report_bytez)
	resp, err := c.DoPost("import-scan", bytez, header)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e *Engagement
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&e); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err)
	}

	return e, nil
}

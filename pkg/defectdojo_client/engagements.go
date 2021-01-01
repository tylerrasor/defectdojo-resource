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

	resp, err := c.DoPost("engagements", engagement_req)
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

type Scan struct {
	Type         string `json:"scan_type"`
	EngagementId int    `json:"engagement"`
	Active       bool   `json:"active"`
}

func (c *DefectdojoClient) UploadReport(report_type string, path string, engagement_id int) error {
	scan_req := Scan{
		Type:         report_type,
		EngagementId: engagement_id,
	}

	resp, err := c.DoPost("import-scan", scan_req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var e *Engagement
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&e); err != nil {
		return fmt.Errorf("error decoding response: %s", err)
	}

	return fmt.Errorf("not implemented")
}

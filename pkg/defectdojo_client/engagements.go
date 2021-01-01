package defectdojo_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
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
	url := fmt.Sprintf("%s/api/v2/engagements/", c.url)
	logrus.Debugf("POST %s", url)

	engagement_req := Engagement{
		ProductId:      p.Id,
		StartDate:      "2021-01-01",
		EndDate:        "2021-01-01",
		EngagementType: "CI/CD",
		EngagementName: report_type,
	}
	bytez, err := json.Marshal(engagement_req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal to json: %s", err)
	}

	logrus.Debugf("trying to send payload: %s", string(bytez))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bytez))
	if err != nil {
		return nil, fmt.Errorf("something went wrong building request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	logrus.Debugln("sending post")
	resp, err := c.DoRequest(req)
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

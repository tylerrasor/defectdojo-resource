package resource

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

type Out struct {
	stdin  io.Reader
	stderr io.Writer
	stdout io.Writer
	args   []string
}

func NewOut(
	stdin io.Reader,
	stderr io.Writer,
	stdout io.Writer,
	args []string,
) *Out {
	return &Out{
		stdin:  stdin,
		stderr: stderr,
		stdout: stdout,
		args:   args,
	}
}

type PutParams struct {
	ReportType string `json:"report_type"`
}

func (p *PutParams) Validate() error {
	if p.ReportType == "" {
		return fmt.Errorf("Required parameter `report_type` not supplied.")
	}

	implemented, key_exists := SupportedReportTypes[p.ReportType]
	if !key_exists {
		return fmt.Errorf("The specified report type, `%s`, is not a supported by Defectdojo (check that your format matches expected)", p.ReportType)
	}
	if !implemented {
		return fmt.Errorf("The specified report type, `%s`, hasn't been implemented yet (pull requests welcome!)", p.ReportType)
	}

	return nil
}

type PutRequest struct {
	Source Source    `json:"source"`
	Params PutParams `json:"params"`
}

func (o *Out) Execute() error {
	logrus.SetOutput(o.stdout)
	var request PutRequest

	decoder := json.NewDecoder(o.stdin)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	if request.Source.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugln("debug logging on")
	}

	logrus.Debugln("getting ready to validate source")
	if err := request.Source.Validate(); err != nil {
		return fmt.Errorf("invalid source config: %s", err)
	}
	logrus.Debugln("source must have validated correctly")

	if err := request.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params config: %s", err)
	}

	return nil
}

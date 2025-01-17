// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/third_party/bazel/gen/bes"
)

type BuildStep interface {
	IsBuildStep()
	GetID() string
	GetStepLabel() string
	GetBuildStepStatus() BuildStepStatus
}

type Problem interface {
	IsNode()
	IsProblem()
	GetID() string
	GetLabel() string
}

type ActionProblem struct {
	ID     string         `json:"id"`
	Label  string         `json:"label"`
	Type   string         `json:"type"`
	Stdout *BlobReference `json:"stdout,omitempty"`
	Stderr *BlobReference `json:"stderr,omitempty"`
	// The underlying BazelInvocationProblem row
	Problem *ent.BazelInvocationProblem `json:"-"`
}

func (ActionProblem) IsNode() {}

func (ActionProblem) IsProblem()            {}
func (this ActionProblem) GetID() string    { return this.ID }
func (this ActionProblem) GetLabel() string { return this.Label }

type BazelCommand struct {
	ID         string `json:"id"`
	Command    string `json:"command"`
	Executable string `json:"executable"`
	Options    string `json:"options"`
	Residual   string `json:"residual"`
}

type BazelInvocationState struct {
	ID             string    `json:"id"`
	BuildEndTime   time.Time `json:"buildEndTime"`
	BuildStartTime time.Time `json:"buildStartTime"`
	ExitCode       *ExitCode `json:"exitCode,omitempty"`
	BepCompleted   bool      `json:"bepCompleted"`
}

type BlobReference struct {
	Name               string             `json:"name"`
	DownloadURL        string             `json:"downloadURL"`
	SizeInBytes        *int               `json:"sizeInBytes,omitempty"`
	AvailabilityStatus ActionOutputStatus `json:"availabilityStatus"`
	// The blob being referenced
	Blob *ent.Blob `json:"-"`
}

type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ExitCode struct {
	ID   string `json:"id"`
	Code int    `json:"code"`
	Name string `json:"name"`
}

type NamedFile struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ProgressProblem struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Output string `json:"output"`
}

func (ProgressProblem) IsNode() {}

func (ProgressProblem) IsProblem()            {}
func (this ProgressProblem) GetID() string    { return this.ID }
func (this ProgressProblem) GetLabel() string { return this.Label }

type TargetProblem struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

func (TargetProblem) IsNode() {}

func (TargetProblem) IsProblem()            {}
func (this TargetProblem) GetID() string    { return this.ID }
func (this TargetProblem) GetLabel() string { return this.Label }

type TestProblem struct {
	ID      string        `json:"id"`
	Label   string        `json:"label"`
	Status  string        `json:"status"`
	Results []*TestResult `json:"results"`
}

func (TestProblem) IsNode() {}

func (TestProblem) IsProblem()            {}
func (this TestProblem) GetID() string    { return this.ID }
func (this TestProblem) GetLabel() string { return this.Label }

type TestResult struct {
	ID                    string         `json:"id"`
	Run                   int            `json:"run"`
	Shard                 int            `json:"shard"`
	Attempt               int            `json:"attempt"`
	Status                string         `json:"status"`
	ActionLogOutput       *BlobReference `json:"actionLogOutput"`
	UndeclaredTestOutputs *BlobReference `json:"undeclaredTestOutputs,omitempty"`
	// TestResult object from the Build Event Stream
	BESTestResult *bes.TestResult `json:"-"`
	// IDs extracted for later use
	TestResultID TestResultID `json:"-"`
}

func (TestResult) IsNode() {}

type User struct {
	ID    string `json:"id"`
	Email string `json:"Email"`
	Ldap  string `json:"LDAP"`
}

type ActionOutputStatus string

const (
	ActionOutputStatusProcessing  ActionOutputStatus = "PROCESSING"
	ActionOutputStatusAvailable   ActionOutputStatus = "AVAILABLE"
	ActionOutputStatusUnavailable ActionOutputStatus = "UNAVAILABLE"
)

var AllActionOutputStatus = []ActionOutputStatus{
	ActionOutputStatusProcessing,
	ActionOutputStatusAvailable,
	ActionOutputStatusUnavailable,
}

func (e ActionOutputStatus) IsValid() bool {
	switch e {
	case ActionOutputStatusProcessing, ActionOutputStatusAvailable, ActionOutputStatusUnavailable:
		return true
	}
	return false
}

func (e ActionOutputStatus) String() string {
	return string(e)
}

func (e *ActionOutputStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ActionOutputStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ActionOutputStatus", str)
	}
	return nil
}

func (e ActionOutputStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type BuildStepStatus string

const (
	BuildStepStatusSuccessful BuildStepStatus = "Successful"
	BuildStepStatusFailed     BuildStepStatus = "Failed"
	BuildStepStatusCancelled  BuildStepStatus = "Cancelled"
	BuildStepStatusUnknown    BuildStepStatus = "Unknown"
)

var AllBuildStepStatus = []BuildStepStatus{
	BuildStepStatusSuccessful,
	BuildStepStatusFailed,
	BuildStepStatusCancelled,
	BuildStepStatusUnknown,
}

func (e BuildStepStatus) IsValid() bool {
	switch e {
	case BuildStepStatusSuccessful, BuildStepStatusFailed, BuildStepStatusCancelled, BuildStepStatusUnknown:
		return true
	}
	return false
}

func (e BuildStepStatus) String() string {
	return string(e)
}

func (e *BuildStepStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = BuildStepStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid BuildStepStatus", str)
	}
	return nil
}

func (e BuildStepStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cron_pub.proto

package pb

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Cron) Validate() error {
	if this.TimeCreated != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.TimeCreated); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("TimeCreated", err)
		}
	}
	if this.TimeUpdated != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.TimeUpdated); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("TimeUpdated", err)
		}
	}
	if this.CronStartTime != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.CronStartTime); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("CronStartTime", err)
		}
	}
	if this.Enable != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Enable); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Enable", err)
		}
	}
	// Validation of proto3 map<> fields is unsupported.
	if this.Extra != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Extra); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Extra", err)
		}
	}
	if this.IsEdge != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.IsEdge); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("IsEdge", err)
		}
	}
	return nil
}
func (this *CronExtra) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	// Validation of proto3 map<> fields is unsupported.
	// Validation of proto3 map<> fields is unsupported.
	// Validation of proto3 map<> fields is unsupported.
	if this.CronStartFrom != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.CronStartFrom); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("CronStartFrom", err)
		}
	}
	if this.Compensator != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Compensator); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Compensator", err)
		}
	}
	if this.LastCompensateAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.LastCompensateAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("LastCompensateAt", err)
		}
	}
	return nil
}
func (this *CronCompensator) Validate() error {
	if this.Enable != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Enable); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Enable", err)
		}
	}
	if this.LatestFirst != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.LatestFirst); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("LatestFirst", err)
		}
	}
	if this.StopIfLatterExecuted != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.StopIfLatterExecuted); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("StopIfLatterExecuted", err)
		}
	}
	return nil
}

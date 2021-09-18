// Code generated by protoc-gen-go-json. DO NOT EDIT.
// Source: build.proto

package pb

import (
	bytes "bytes"
	json "encoding/json"

	jsonpb "github.com/erda-project/erda-infra/pkg/transport/http/encoding/jsonpb"
	protojson "google.golang.org/protobuf/encoding/protojson"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the "encoding/json" package it is being compiled against.
var _ json.Marshaler = (*BuildCacheReportRequest)(nil)
var _ json.Unmarshaler = (*BuildCacheReportRequest)(nil)
var _ json.Marshaler = (*BuildCacheReportResponse)(nil)
var _ json.Unmarshaler = (*BuildCacheReportResponse)(nil)
var _ json.Marshaler = (*BuildArtifactRegisterRequest)(nil)
var _ json.Unmarshaler = (*BuildArtifactRegisterRequest)(nil)
var _ json.Marshaler = (*BuildArtifactRegisterResponse)(nil)
var _ json.Unmarshaler = (*BuildArtifactRegisterResponse)(nil)
var _ json.Marshaler = (*BuildArtifactQueryRequest)(nil)
var _ json.Unmarshaler = (*BuildArtifactQueryRequest)(nil)
var _ json.Marshaler = (*BuildArtifactQueryResponse)(nil)
var _ json.Unmarshaler = (*BuildArtifactQueryResponse)(nil)
var _ json.Marshaler = (*BuildArtifact)(nil)
var _ json.Unmarshaler = (*BuildArtifact)(nil)

// BuildCacheReportRequest implement json.Marshaler.
func (m *BuildCacheReportRequest) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := (&jsonpb.Marshaler{
		OrigName:     false,
		EnumsAsInts:  false,
		EmitDefaults: true,
	}).Marshal(buf, m)
	return buf.Bytes(), err
}

// BuildCacheReportRequest implement json.Marshaler.
func (m *BuildCacheReportRequest) UnmarshalJSON(b []byte) error {
	return (&protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}).Unmarshal(b, m)
}

// BuildCacheReportResponse implement json.Marshaler.
func (m *BuildCacheReportResponse) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := (&jsonpb.Marshaler{
		OrigName:     false,
		EnumsAsInts:  false,
		EmitDefaults: true,
	}).Marshal(buf, m)
	return buf.Bytes(), err
}

// BuildCacheReportResponse implement json.Marshaler.
func (m *BuildCacheReportResponse) UnmarshalJSON(b []byte) error {
	return (&protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}).Unmarshal(b, m)
}

// BuildArtifactRegisterRequest implement json.Marshaler.
func (m *BuildArtifactRegisterRequest) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := (&jsonpb.Marshaler{
		OrigName:     false,
		EnumsAsInts:  false,
		EmitDefaults: true,
	}).Marshal(buf, m)
	return buf.Bytes(), err
}

// BuildArtifactRegisterRequest implement json.Marshaler.
func (m *BuildArtifactRegisterRequest) UnmarshalJSON(b []byte) error {
	return (&protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}).Unmarshal(b, m)
}

// BuildArtifactRegisterResponse implement json.Marshaler.
func (m *BuildArtifactRegisterResponse) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := (&jsonpb.Marshaler{
		OrigName:     false,
		EnumsAsInts:  false,
		EmitDefaults: true,
	}).Marshal(buf, m)
	return buf.Bytes(), err
}

// BuildArtifactRegisterResponse implement json.Marshaler.
func (m *BuildArtifactRegisterResponse) UnmarshalJSON(b []byte) error {
	return (&protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}).Unmarshal(b, m)
}

// BuildArtifactQueryRequest implement json.Marshaler.
func (m *BuildArtifactQueryRequest) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := (&jsonpb.Marshaler{
		OrigName:     false,
		EnumsAsInts:  false,
		EmitDefaults: true,
	}).Marshal(buf, m)
	return buf.Bytes(), err
}

// BuildArtifactQueryRequest implement json.Marshaler.
func (m *BuildArtifactQueryRequest) UnmarshalJSON(b []byte) error {
	return (&protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}).Unmarshal(b, m)
}

// BuildArtifactQueryResponse implement json.Marshaler.
func (m *BuildArtifactQueryResponse) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := (&jsonpb.Marshaler{
		OrigName:     false,
		EnumsAsInts:  false,
		EmitDefaults: true,
	}).Marshal(buf, m)
	return buf.Bytes(), err
}

// BuildArtifactQueryResponse implement json.Marshaler.
func (m *BuildArtifactQueryResponse) UnmarshalJSON(b []byte) error {
	return (&protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}).Unmarshal(b, m)
}

// BuildArtifact implement json.Marshaler.
func (m *BuildArtifact) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := (&jsonpb.Marshaler{
		OrigName:     false,
		EnumsAsInts:  false,
		EmitDefaults: true,
	}).Marshal(buf, m)
	return buf.Bytes(), err
}

// BuildArtifact implement json.Marshaler.
func (m *BuildArtifact) UnmarshalJSON(b []byte) error {
	return (&protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}).Unmarshal(b, m)
}

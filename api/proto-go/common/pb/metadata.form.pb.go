// Code generated by protoc-gen-go-form. DO NOT EDIT.
// Source: metadata.proto

package pb

import (
	url "net/url"
	strconv "strconv"

	urlenc "github.com/erda-project/erda-infra/pkg/urlenc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the "github.com/erda-project/erda-infra/pkg/urlenc" package it is being compiled against.
var _ urlenc.URLValuesUnmarshaler = (*MetadataField)(nil)

// MetadataField implement urlenc.URLValuesUnmarshaler.
func (m *MetadataField) UnmarshalURLValues(prefix string, values url.Values) error {
	for key, vals := range values {
		if len(vals) > 0 {
			switch prefix + key {
			case "name":
				m.Name = vals[0]
			case "value":
				m.Value = vals[0]
			case "type":
				m.Type = vals[0]
			case "optional":
				val, err := strconv.ParseBool(vals[0])
				if err != nil {
					return err
				}
				m.Optional = val
			case "level":
				m.Level = vals[0]
			}
		}
	}
	return nil
}

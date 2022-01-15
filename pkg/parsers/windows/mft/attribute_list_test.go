// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func Test_getAttributeListAttribute(t *testing.T) {
	tests := []struct {
		name                        string
		input                       []byte
		wantAttributeListAttributes AttributeListAttributes
		wantErr                     bool
	}{
		{
			name:    "test1",
			input:   []byte{0x20, 0x00, 0x00, 0x00, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x00, 0x80, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0x87, 0x07, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x44, 0x43, 0x00, 0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0x87, 0x07, 0x00, 0x00, 0x00, 0x08, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0xB7, 0x15, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x28, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x45, 0xB7, 0x15, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			wantErr: false,
			wantAttributeListAttributes: AttributeListAttributes{
				AttributeListAttribute{
					Type:                     0x10,
					MFTReferenceRecordNumber: 493401,
				},
				AttributeListAttribute{
					Type:                     0x30,
					MFTReferenceRecordNumber: 493401,
				},
				AttributeListAttribute{
					Type:                     0x80,
					MFTReferenceRecordNumber: 1423172,
				},
				AttributeListAttribute{
					Type:                     0x80,
					MFTReferenceRecordNumber: 1423173,
				},
			},
		},
		{
			name:                        "nil []byte test",
			input:                       nil,
			wantErr:                     true,
			wantAttributeListAttributes: AttributeListAttributes{},
		},
		{
			name:                        "wrong attribute test",
			input:                       []byte{0x10, 0x00, 0x00, 0x00, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x00, 0x80, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0x87, 0x07, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x44, 0x43, 0x00, 0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0x87, 0x07, 0x00, 0x00, 0x00, 0x08, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0xB7, 0x15, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x28, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x45, 0xB7, 0x15, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			wantErr:                     true,
			wantAttributeListAttributes: AttributeListAttributes{},
		},
		{
			name:                        "unexpected size",
			input:                       []byte{0x20, 0x00, 0x00, 0x00, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x00, 0x80, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0x87, 0x07, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x44, 0x43, 0x00, 0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0x87, 0x07, 0x00, 0x00, 0x00, 0x08, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0xB7, 0x15, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x28, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x45, 0xB7, 0x15, 0x00, 0x00},
			wantErr:                     true,
			wantAttributeListAttributes: AttributeListAttributes{},
		},
		{
			name:    "unexpected subattribute",
			input:   []byte{0x20, 0x00, 0x00, 0x00, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x00, 0x80, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0x87, 0x07, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x44, 0x43, 0x00, 0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0x87, 0x07, 0x00, 0x00, 0x00, 0x08, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0xB7, 0x15, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x81, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x28, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x45, 0xB7, 0x15, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			wantErr: false,
			wantAttributeListAttributes: AttributeListAttributes{
				AttributeListAttribute{
					Type:                     0x10,
					MFTReferenceRecordNumber: 493401,
				},
				AttributeListAttribute{
					Type:                     0x30,
					MFTReferenceRecordNumber: 493401,
				},
				AttributeListAttribute{
					Type:                     0x80,
					MFTReferenceRecordNumber: 1423172,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAttributeListAttributes, err := getAttributeListAttribute(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAttributeListAttributes, tt.wantAttributeListAttributes) {
				t.Errorf(cmp.Diff(gotAttributeListAttributes, tt.wantAttributeListAttributes))
			}
		})
	}
}

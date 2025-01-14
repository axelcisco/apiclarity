// Copyright © 2021 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"io"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/openclarity/apiclarity/plugins/api/client/models"
)

func Test_GetPathWithQuery(t *testing.T) {
	type args struct {
		reqURL *url.URL
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no query",
			args: args{
				reqURL: &url.URL{
					Path:     "/foo/bar",
					RawQuery: "",
				},
			},
			want: "/foo/bar",
		},
		{
			name: "with query",
			args: args{
				reqURL: &url.URL{
					Path:     "/foo/bar",
					RawQuery: "bla=bloo",
				},
			},
			want: "/foo/bar?bla=bloo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPathWithQuery(tt.args.reqURL); got != tt.want {
				t.Errorf("getPathWithQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ReadBody(t *testing.T) {
	reqBodyJSON := "{Hello: world!}"
	reqBodyJSONLong := "aaaaaaaaaaaaaaaaaaaa"
	for i := 0; i < 16; i++ {
		reqBodyJSONLong += reqBodyJSONLong
	}
	type args struct {
		body io.ReadCloser
	}
	tests := []struct {
		name          string
		args          args
		want          []byte
		wantTruncated bool
		wantErr       bool
	}{
		{
			name: "body is not truncated",
			args: args{
				body: io.NopCloser(strings.NewReader(reqBodyJSON)),
			},
			want:          []byte(reqBodyJSON),
			wantTruncated: false,
			wantErr:       false,
		},
		{
			name: "body is truncated",
			args: args{
				body: io.NopCloser(strings.NewReader(reqBodyJSONLong)),
			},
			want:          []byte{},
			wantTruncated: true,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ReadBody(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("readBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readBody() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantTruncated {
				t.Errorf("readBody() got1 = %v, want %v", got1, tt.wantTruncated)
			}
		})
	}
}

func TestCreateHeaders(t *testing.T) {
	type args struct {
		headers map[string][]string
	}
	tests := []struct {
		name string
		args args
		want []*models.Header
	}{
		{
			name: "",
			args: args{
				headers: map[string][]string{
					"h1": {"v1"},
					"h2": {"v2"},
				},
			},
			want: []*models.Header{
				{
					Key:   "h1",
					Value: "v1",
				},
				{
					Key:   "h2",
					Value: "v2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateHeaders(tt.args.headers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

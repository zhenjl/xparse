// Copyright (c) 2014 Dataence, LLC. All rights reserved.
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

package xtld

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/publicsuffix"
)

func TestXtldTldlen(t *testing.T) {
	testdata := map[string]string{
		"abc.com":       "com",
		"www.abc.com":   "com",
		"中国":            "中国",
		"england.co.uk": "co.uk",
		"xxxxx":         "",
	}

	for k, v := range testdata {
		sfx := TLD(k)
		require.Equal(t, sfx, v)
		require.Equal(t, len(sfx), len(v))

		// PublicSuffix fails for "xxxxx"
		//psfx, _ := publicsuffix.PublicSuffix(k)
		//require.Equal(t, sfx, psfx)
	}
}

func BenchmarkXtldTldLen(b *testing.B) {
	//ss := "england.co.uk"
	ss := "中国"
	for i := 0; i < b.N; i++ {
		TLDlen(ss)
	}
}

func BenchmarkXtldTld(b *testing.B) {
	//ss := "england.co.uk"
	ss := "中国"
	for i := 0; i < b.N; i++ {
		TLD(ss)
	}
}

func BenchmarkPublicsuffix(b *testing.B) {
	ss := "中国"
	for i := 0; i < b.N; i++ {
		publicsuffix.PublicSuffix(ss)
	}
}

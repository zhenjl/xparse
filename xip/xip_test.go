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

package netx

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/surge/glog"
)

func TestParseIPv4Success(t *testing.T) {
	for i, ip := range ips {
		glog.Debugf("Parsing %d %s", i, ip)
		res, err := Parse(ip)
		require.NoError(t, err)

		m := make(map[[4]byte]bool)

		for _, ip2 := range res {
			var tmp [4]byte
			copy(tmp[:], ip2.To4())
			m[tmp] = true
		}

		require.Equal(t, results[i], m)
	}
}

func TestParseIPv4Failure(t *testing.T) {
	_, err := Parse("10.1.1.1,10.1.1.2")
	require.Error(t, err)

	_, err = Parse("10.1.1.1.")
	require.Error(t, err)

	_, err = Parse("10.1.1.a")
	require.Error(t, err)

	_, err = Parse("10.1.1.256")
	require.Error(t, err)
}

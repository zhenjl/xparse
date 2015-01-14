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

package timex

import (
	"regexp"
	"testing"
	"time"

	"github.com/dataence/assert"
)

var (
	re1 = regexp.MustCompile("_")
	re2 = regexp.MustCompile("Z")
)

func TestTimeFormats(t *testing.T) {
	for _, f := range TimeFormats {
		tx := re2.ReplaceAllString(re1.ReplaceAllString(f, " "), "+")
		expected, err := time.Parse(f, tx)
		assert.NoError(t, true, err)
		actual, err := Parse(tx)
		assert.NoError(t, true, err)
		assert.Equal(t, true, expected.UnixNano(), actual.UnixNano())
	}
}

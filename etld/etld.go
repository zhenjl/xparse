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

// etld is an effective TLD matcher that returns the length of the effective domain
// name for the given string. It uses the data set from
// https://www.publicsuffix.org/list/effective_tld_names.dat
//
// For extracting public suffxies, please use golang.org/x/net/publicsuffix.
package etld

import (
	"strings"

	"golang.org/x/net/idna"
)

//go:generate go run genfsm.go effective_tld_names.dat.gz fsm.go

// Match returns the length of the ETLD matched, if no match, then 0 is returned.
func Match(s string) int {
	s, err := idna.ToUnicode(strings.ToLower(s))
	if err != nil {
		return 0
	}

	return etldlen(s)
}

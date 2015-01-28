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

// etld is an effective TLD parser that extracts the top-level-domain out from
// the given string. It uses the data set from
// https://www.publicsuffix.org/list/effective_tld_names.dat
//
// Unfortunately it's still not 100% and there's already a working implementation
// at golang.org/x/net/publicsuffix. So I recommend using that instead.
package etld

import (
	"fmt"
	"strings"

	"golang.org/x/net/idna"
)

//go:generate go run genfsm.go effective_tld_names.dat.gz fsm.go

// RegistrableDomain returns the the registered or registrable domain, which is the
// public suffix plus one additional label. If the domain matches an ETLD, then
// the registrable domain is "".
func RegistrableDomain(ss string) (string, error) {
	if ss == "" {
		return "", nil
	}

	ss, err := idna.ToUnicode(strings.ToLower(ss))
	if err != nil {
		return "", err
	}

	l := len(ss)
	m := etldlen(ss)
	fmt.Printf("ss=%q, l=%d, m=%d\n", ss, l, m)

	if len(ss) == m {
		return "", nil
	} else if len(ss)-m == 1 && ss[0] == '.' {
		return "", fmt.Errorf("etld: invalid domain %s", ss)
	}

	i := l - m - 1
	if ss[i] != '.' {
		return "", fmt.Errorf("etld: invalid suffix %s for domain %s", ss[l-m:], ss)
	}

	fmt.Printf("ss=%q, l=%d, m=%d\n", ss[:i], l, m)

	return idna.ToASCII(ss[1+strings.LastIndex(ss[:i], "."):])
}

// PublicSuffix returns the public suffix of a domain based on the algorithm from
// https://publicsuffix.org/list/. This may return something different from ETLD().
// Namely, anything that returns "" by ETLD() will return the original string by
// PublicSuffix(). Review algorithm for more details.
func PublicSuffix(ss string) (string, error) {
	ss, err := idna.ToUnicode(strings.ToLower(ss))
	if err != nil {
		return "", err
	}

	m := etldlen(ss)
	// If no rules match, the prevailing rule is "*".
	if m == 0 {
		return ss, nil
	}
	return ss[len(ss)-m:], nil
}

// ETLD extracts the top-level-domain based on data from
// https://www.publicsuffix.org/list/effective_tld_names.dat
func ETLD(ss string) (string, error) {
	ss, err := idna.ToUnicode(strings.ToLower(ss))
	if err != nil {
		return "", err
	}

	return etld(ss), nil
}

// etldlen determins the length of the TLD
func ETLDlen(ss string) (int, error) {
	ss, err := idna.ToUnicode(strings.ToLower(ss))
	if err != nil {
		return 0, err
	}

	return etldlen(ss), nil
}

func etld(ss string) string {
	m := etldlen(ss)
	return ss[len(ss)-m:]
}

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

package etld

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
)

// Copied from https://github.com/golang/net/blob/master/publicsuffix/list_test.go
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found at https://github.com/golang/net/blob/master/LICENSE
var publicSuffixTestCases = []struct {
	domain, want string
}{
	// Empty string.
	{"", ""},

	// The .ao rules are:
	// ao
	// ed.ao
	// gv.ao
	// og.ao
	// co.ao
	// pb.ao
	// it.ao
	{"ao", "ao"},
	{"www.ao", "ao"},
	{"pb.ao", "pb.ao"},
	{"www.pb.ao", "pb.ao"},
	{"www.xxx.yyy.zzz.pb.ao", "pb.ao"},

	// The .ar rules are:
	// ar
	// com.ar
	// edu.ar
	// gob.ar
	// gov.ar
	// int.ar
	// mil.ar
	// net.ar
	// org.ar
	// tur.ar
	// blogspot.com.ar
	{"ar", "ar"},
	{"www.ar", "ar"},
	{"nic.ar", "ar"},
	{"www.nic.ar", "ar"},
	{"com.ar", "com.ar"},
	{"www.com.ar", "com.ar"},
	{"blogspot.com.ar", "blogspot.com.ar"},
	{"www.blogspot.com.ar", "blogspot.com.ar"},
	{"www.xxx.yyy.zzz.blogspot.com.ar", "blogspot.com.ar"},
	{"logspot.com.ar", "com.ar"},
	{"zlogspot.com.ar", "com.ar"},
	{"zblogspot.com.ar", "com.ar"},

	// The .arpa rules are:
	// arpa
	// e164.arpa
	// in-addr.arpa
	// ip6.arpa
	// iris.arpa
	// uri.arpa
	// urn.arpa
	{"arpa", "arpa"},
	{"www.arpa", "arpa"},
	{"urn.arpa", "urn.arpa"},
	{"www.urn.arpa", "urn.arpa"},
	{"www.xxx.yyy.zzz.urn.arpa", "urn.arpa"},

	// The relevant {kobe,kyoto}.jp rules are:
	// jp
	// *.kobe.jp
	// !city.kobe.jp
	// kyoto.jp
	// ide.kyoto.jp
	{"jp", "jp"},
	{"kobe.jp", "jp"},
	{"c.kobe.jp", "c.kobe.jp"},
	{"b.c.kobe.jp", "c.kobe.jp"},
	{"a.b.c.kobe.jp", "c.kobe.jp"},
	{"city.kobe.jp", "kobe.jp"},
	{"www.city.kobe.jp", "kobe.jp"},
	{"kyoto.jp", "kyoto.jp"},
	{"test.kyoto.jp", "kyoto.jp"},
	{"ide.kyoto.jp", "ide.kyoto.jp"},
	{"b.ide.kyoto.jp", "ide.kyoto.jp"},
	{"a.b.ide.kyoto.jp", "ide.kyoto.jp"},

	// The .tw rules are:
	// tw
	// edu.tw
	// gov.tw
	// mil.tw
	// com.tw
	// net.tw
	// org.tw
	// idv.tw
	// game.tw
	// ebiz.tw
	// club.tw
	// 網路.tw (xn--zf0ao64a.tw)
	// 組織.tw (xn--uc0atv.tw)
	// 商業.tw (xn--czrw28b.tw)
	// blogspot.tw
	{"tw", "tw"},
	{"aaa.tw", "tw"},
	{"www.aaa.tw", "tw"},
	//{"xn--czrw28b.aaa.tw", "tw"},
	{"edu.tw", "edu.tw"},
	{"www.edu.tw", "edu.tw"},
	//{"xn--czrw28b.edu.tw", "edu.tw"},
	//{"xn--czrw28b.tw", "xn--czrw28b.tw"},
	//{"www.xn--czrw28b.tw", "xn--czrw28b.tw"},
	//{"xn--uc0atv.xn--czrw28b.tw", "xn--czrw28b.tw"},
	//{"xn--kpry57d.tw", "tw"},

	// The .uk rules are:
	// uk
	// ac.uk
	// co.uk
	// gov.uk
	// ltd.uk
	// me.uk
	// net.uk
	// nhs.uk
	// org.uk
	// plc.uk
	// police.uk
	// *.sch.uk
	// blogspot.co.uk
	{"uk", "uk"},
	{"aaa.uk", "uk"},
	{"www.aaa.uk", "uk"},
	{"mod.uk", "uk"},
	{"www.mod.uk", "uk"},
	{"sch.uk", "uk"},
	{"mod.sch.uk", "mod.sch.uk"},
	{"www.sch.uk", "www.sch.uk"},
	{"blogspot.co.uk", "blogspot.co.uk"},
	{"blogspot.nic.uk", "uk"},
	{"blogspot.sch.uk", "blogspot.sch.uk"},

	// The .рф rules are
	// рф (xn--p1ai)
	//{"xn--p1ai", "xn--p1ai"},
	//{"aaa.xn--p1ai", "xn--p1ai"},
	//{"www.xxx.yyy.xn--p1ai", "xn--p1ai"},

	// The .zw rules are:
	// *.zw
	{"zw", "zw"},
	{"www.zw", "www.zw"},
	{"zzz.zw", "zzz.zw"},
	{"www.zzz.zw", "zzz.zw"},
	{"www.xxx.yyy.zzz.zw", "zzz.zw"},

	// There are no .nosuchtld rules.
	{"nosuchtld", "nosuchtld"},
	{"foo.nosuchtld", "nosuchtld"},
	{"bar.foo.nosuchtld", "nosuchtld"},
}

var registrableDomainTestCases = []struct {
	domain, want string
}{
	// The following list is from
	// http://mxr.mozilla.org/mozilla-central/source/netwerk/test/unit/data/test_psl.txt?raw=1
	// http://creativecommons.org/publicdomain/zero/1.0/
	// "" input.
	{"", ""},
	// Mixed case.
	//{"COM", ""},
	//{"example.COM", "example.com"},
	//{"WwW.example.COM", "example.com"},
	// Leading dot.
	//{".com", ""},
	//{".example", ""},
	//{".example.com", ""},
	//{".example.example", ""},
	// Unlisted TLD.
	{"example", ""},
	{"example.example", "example.example"},
	{"b.example.example", "example.example"},
	{"a.b.example.example", "example.example"},
	// Listed, but non-Internet, TLD.
	//"local", ""},
	//"example.local", ""},
	//"b.example.local", ""},
	//"a.b.example.local", ""},
	// TLD with only 1 rule.
	{"biz", ""},
	{"domain.biz", "domain.biz"},
	{"b.domain.biz", "domain.biz"},
	{"a.b.domain.biz", "domain.biz"},
	// TLD with some 2-level rules.
	{"com", ""},
	{"example.com", "example.com"},
	{"b.example.com", "example.com"},
	{"a.b.example.com", "example.com"},
	{"uk.com", ""},
	{"example.uk.com", "example.uk.com"},
	{"b.example.uk.com", "example.uk.com"},
	{"a.b.example.uk.com", "example.uk.com"},
	{"test.ac", "test.ac"},
	// TLD with only 1 (wildcard, rule.
	{"cy", ""},
	{"c.cy", ""},
	{"b.c.cy", "b.c.cy"},
	{"a.b.c.cy", "b.c.cy"},
	// More complex TLD.
	{"jp", ""},
	{"test.jp", "test.jp"},
	{"www.test.jp", "test.jp"},
	{"ac.jp", ""},
	{"test.ac.jp", "test.ac.jp"},
	{"www.test.ac.jp", "test.ac.jp"},
	{"kyoto.jp", ""},
	{"test.kyoto.jp", "test.kyoto.jp"},
	{"ide.kyoto.jp", ""},
	{"b.ide.kyoto.jp", "b.ide.kyoto.jp"},
	{"a.b.ide.kyoto.jp", "b.ide.kyoto.jp"},
	{"c.kobe.jp", ""},
	{"b.c.kobe.jp", "b.c.kobe.jp"},
	{"a.b.c.kobe.jp", "b.c.kobe.jp"},
	{"city.kobe.jp", "city.kobe.jp"},
	{"www.city.kobe.jp", "city.kobe.jp"},
	// TLD with a wildcard rule and exceptions.
	{"ck", ""},
	{"test.ck", ""},
	{"b.test.ck", "b.test.ck"},
	{"a.b.test.ck", "b.test.ck"},
	{"www.ck", "www.ck"},
	{"www.www.ck", "www.ck"},
	// US K12.
	{"us", ""},
	{"test.us", "test.us"},
	{"www.test.us", "test.us"},
	{"ak.us", ""},
	{"test.ak.us", "test.ak.us"},
	{"www.test.ak.us", "test.ak.us"},
	{"k12.ak.us", ""},
	{"test.k12.ak.us", "test.k12.ak.us"},
	{"www.test.k12.ak.us", "test.k12.ak.us"},
	// IDN labels.
	{"食狮.com.cn", "食狮.com.cn"},
	//{"食狮.公司.cn", "食狮.公司.cn"},
	//{"www.食狮.公司.cn", "食狮.公司.cn"},
	//{"shishi.公司.cn", "shishi.公司.cn"},
	//{"公司.cn", ""},
	{"食狮.中国", "食狮.中国"},
	{"www.食狮.中国", "食狮.中国"},
	{"shishi.中国", "shishi.中国"},
	{"中国", ""},
	// Same as above, but punycoded. Not supported.
	{"xn--85x722f.com.cn", "xn--85x722f.com.cn"},
	{"xn--85x722f.xn--55qx5d.cn", "xn--85x722f.xn--55qx5d.cn"},
	{"www.xn--85x722f.xn--55qx5d.cn", "xn--85x722f.xn--55qx5d.cn"},
	{"shishi.xn--55qx5d.cn", "shishi.xn--55qx5d.cn"},
	{"xn--55qx5d.cn", ""},
	{"xn--85x722f.xn--fiqs8s", "xn--85x722f.xn--fiqs8s"},
	{"www.xn--85x722f.xn--fiqs8s", "xn--85x722f.xn--fiqs8s"},
	{"shishi.xn--fiqs8s", "shishi.xn--fiqs8s"},
	{"xn--fiqs8s", ""},
}

func TestEtldPublicSuffix(t *testing.T) {
	for _, s := range publicSuffixTestCases {
		psfx, err := PublicSuffix(s.domain)
		require.NoError(t, err, s.domain)
		require.Equal(t, s.want, psfx, s.domain)
	}
}

func TestEtldRegistrableDomain(t *testing.T) {
	for _, s := range registrableDomainTestCases {
		dom, _ := RegistrableDomain(s.domain)
		//require.NoError(t, err, s.domain)
		require.Equal(t, s.want, dom, s.domain)
	}
}

func TestXnetRegistrableDomain(t *testing.T) {
	for _, s := range registrableDomainTestCases {
		dom, _ := publicsuffix.EffectiveTLDPlusOne(s.domain)
		//require.NoError(t, err, s.domain)
		require.Equal(t, s.want, dom, s.domain)
	}
}

func BenchmarkXparsePublicSuffix(b *testing.B) {
	ss := "aaaaaaaaaaaaaaaabc.england.co.uk"
	//ss := "中国"
	for i := 0; i < b.N; i++ {
		PublicSuffix(ss)
	}
}

func BenchmarkXnetPublicSuffix(b *testing.B) {
	ss := "aaaaaaaaaaaaaaaabc.england.co.uk"
	//ss := "中国"
	for i := 0; i < b.N; i++ {
		publicsuffix.PublicSuffix(ss)
	}
}

func TestIdnaToUnicode(t *testing.T) {
	fmt.Println(idna.ToUnicode("shishi.xn--fiqs8s"))
	fmt.Println(idna.ToUnicode("中国"))
	fmt.Println(publicsuffix.PublicSuffix("中国"))
	fmt.Println(publicsuffix.PublicSuffix("shishi.xn--fiqs8s"))
}

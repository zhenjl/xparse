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

// xtld is a TLD parser that extracts the top-level-domain out from the
// given string. It uses the data set from
// https://www.publicsuffix.org/list/effective_tld_names.dat
//
// This package is inspired by https://github.com/jpillora/go-tld.
//
// Example:
//
//   import "github.com/parse/xtld"
//
//   xtld.TLD("abc.com") // should return "com"
//   xtld.TLD("中国") // should return "中国"
//   xtld.TLD("xxxxx") // should return ""
//
// Due to the large number of states, the tldfsm.go file is over 100K lines.
// It takes a bit of time to compile but it's worth it. Benchmark shows that
// each extraction takes only 50-70ns.
package xtld

//go:generate xtldfsm effective_tld_names.dat tldfsm.go

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

// xparse is a collection of parsing utilities, such as IP, time, and top-level-domain, in Go.
//
// - xip is an IP (only v4 for now) parser that expands a given IP string to all the IP addresses it represents.
//
// - xtime is a time parser that parses the time without knowning the exact format.
//
// - xtld is a TLD parser that extracts the top-level-domain out from the given string.
package xparse

// Copyright © 2016 Abcum Ltd
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

package fibre

import (
	"reflect"

	"github.com/ugorji/go/codec"
)

var (
	jh codec.JsonHandle
	ch codec.CborHandle
	bh codec.BincHandle
	mh codec.MsgpackHandle
)

func init() {

	// JSONHandle

	jh.Canonical = true
	jh.InternString = true
	jh.CheckCircularRef = false
	jh.AsSymbols = codec.AsSymbolDefault
	jh.SliceType = reflect.TypeOf([]interface{}(nil))
	jh.MapType = reflect.TypeOf(map[string]interface{}(nil))

	// CBORHandle

	ch.Canonical = true
	ch.InternString = true
	ch.CheckCircularRef = false
	ch.AsSymbols = codec.AsSymbolDefault
	ch.SliceType = reflect.TypeOf([]interface{}(nil))
	ch.MapType = reflect.TypeOf(map[string]interface{}(nil))

	// BINCHandle

	bh.Canonical = true
	bh.InternString = true
	bh.CheckCircularRef = false
	bh.AsSymbols = codec.AsSymbolDefault
	bh.SliceType = reflect.TypeOf([]interface{}(nil))
	bh.MapType = reflect.TypeOf(map[string]interface{}(nil))

	// PACKHandle

	mh.WriteExt = true
	mh.Canonical = true
	mh.RawToString = true
	mh.InternString = true
	mh.CheckCircularRef = false
	mh.AsSymbols = codec.AsSymbolDefault
	mh.SliceType = reflect.TypeOf([]interface{}(nil))
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))

}

//  Copyright (c) 2015 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package firestorm

import (
	"testing"

	"github.com/blevesearch/bleve/index/store"
	"github.com/blevesearch/bleve/index/store/gtreap"
)

func CreateGTreap() (store.KVStore, error) {
	return gtreap.StoreConstructor(nil)
}

func DestroyGTreap() error {
	return nil
}

func BenchmarkGTreapIndexing1Workers(b *testing.B) {
	CommonBenchmarkIndex(b, CreateGTreap, DestroyGTreap, 1)
}

func BenchmarkGTreapIndexing2Workers(b *testing.B) {
	CommonBenchmarkIndex(b, CreateGTreap, DestroyGTreap, 2)
}

func BenchmarkGTreapIndexing4Workers(b *testing.B) {
	CommonBenchmarkIndex(b, CreateGTreap, DestroyGTreap, 4)
}

// batches

func BenchmarkGTreapIndexing1Workers10Batch(b *testing.B) {
	CommonBenchmarkIndexBatch(b, CreateGTreap, DestroyGTreap, 1, 10)
}

func BenchmarkGTreapIndexing2Workers10Batch(b *testing.B) {
	CommonBenchmarkIndexBatch(b, CreateGTreap, DestroyGTreap, 2, 10)
}

func BenchmarkGTreapIndexing4Workers10Batch(b *testing.B) {
	CommonBenchmarkIndexBatch(b, CreateGTreap, DestroyGTreap, 4, 10)
}

func BenchmarkGTreapIndexing1Workers100Batch(b *testing.B) {
	CommonBenchmarkIndexBatch(b, CreateGTreap, DestroyGTreap, 1, 100)
}

func BenchmarkGTreapIndexing2Workers100Batch(b *testing.B) {
	CommonBenchmarkIndexBatch(b, CreateGTreap, DestroyGTreap, 2, 100)
}

func BenchmarkGTreapIndexing4Workers100Batch(b *testing.B) {
	CommonBenchmarkIndexBatch(b, CreateGTreap, DestroyGTreap, 4, 100)
}

func BenchmarkGTreapIndexing1Workers1000Batch(b *testing.B) {
	CommonBenchmarkIndexBatch(b, CreateGTreap, DestroyGTreap, 1, 1000)
}

func BenchmarkGTreapIndexing2Workers1000Batch(b *testing.B) {
	CommonBenchmarkIndexBatch(b, CreateGTreap, DestroyGTreap, 2, 1000)
}

func BenchmarkGTreapIndexing4Workers1000Batch(b *testing.B) {
	CommonBenchmarkIndexBatch(b, CreateGTreap, DestroyGTreap, 4, 1000)
}

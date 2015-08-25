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
	"bytes"
	"fmt"
	"sort"
	"sync/atomic"

	"github.com/blevesearch/bleve/index/store"
)

const IDFieldName = "_id"

func (f *Firestorm) bootstrap(writer store.KVWriter) error {
	// record version
	err := f.storeVersion(writer)
	if err != nil {
		return err
	}
	// define _id field
	_, idFieldRow := f.fieldIndexOrNewRow(IDFieldName)
	err = writer.Set(idFieldRow.Key(), idFieldRow.Value())
	if err != nil {
		return err
	}

	return nil
}

func (f *Firestorm) warmup(writer store.KVWriter) error {
	// load all the existing fields
	err := f.loadFields(writer)
	if err != nil {
		return err
	}

	// walk the term frequency info for _id
	// this allows us to find deleted doc numbers
	// and seed the doc count
	idField, existed := f.fieldCache.FieldNamed(IDFieldName, false)
	if !existed {
		return fmt.Errorf("_id field missing, cannot proceed")
	}

	tfkPrefix := TermFreqIteratorStart(idField, nil)

	var lastDocId []byte
	lastDocNumbers := make(DocNumberList, 1)
	err = visitPrefix(writer, tfkPrefix, func(key, val []byte) (bool, error) {
		tfk, err := NewTermFreqRowKV(key, val)
		if err != nil {
			return false, err
		}
		docID := tfk.DocID()
		docNum := tfk.DocNum()

		if docNum > f.highDocNumber {
			f.highDocNumber = docNum
		}
		if docNum > f.compensator.maxRead {
			f.compensator.maxRead = docNum
		}

		// check for consecutive records
		if bytes.Compare(docID, lastDocId) == 0 {
			lastDocNumbers = append(lastDocNumbers, docNum)
		} else {
			// new doc id
			atomic.AddUint64(f.docCount, 1)

			// last docID had multiple doc numbers
			if len(lastDocNumbers) > 1 {
				f.addOldDocNumbers(lastDocNumbers, lastDocId)

				// reset size to 1
				lastDocNumbers = make(DocNumberList, 1)
			}
			lastDocNumbers[0] = docNum
			lastDocId = docID
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	// be sure to finish up check on final row
	if len(lastDocNumbers) > 1 {
		f.addOldDocNumbers(lastDocNumbers, lastDocId)
	}

	return nil
}

func (f *Firestorm) addOldDocNumbers(docNumberList DocNumberList, docID []byte) {
	sort.Sort(docNumberList)
	// high doc number is OK, rest are deleted
	for _, dn := range docNumberList[1:] {
		// f.deletedDocNumbers.Add(dn, docID)
		f.compensator.deletedDocNumbers.Set(uint(dn))
		f.garbageCollector.Notify(dn, docID)
	}
}

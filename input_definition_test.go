// Copyright 2017 Pilosa Corp.
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

package pilosa_test

import (
	"encoding/json"
	"testing"

	"github.com/pilosa/pilosa"
	"github.com/pilosa/pilosa/internal"
)

func TestInputDefinition_Open(t *testing.T) {
	index := MustOpenIndex()
	defer index.Close()

	// Create Input Definition.
	frames := internal.Frame{Name: "f", Meta: &internal.FrameMeta{RowLabel: "row"}}
	action := internal.Action{Frame: "f", ValueDestination: "map", ValueMap: map[string]uint64{"Green": 1}}
	fields := internal.InputDefinitionField{Name: "id", PrimaryKey: true, Actions: []*internal.Action{&action}}
	def := internal.InputDefinition{Name: "test", Frames: []*internal.Frame{&frames}, Fields: []*internal.InputDefinitionField{&fields}}
	inputDef, err := index.CreateInputDefinition(&def)
	if err != nil {
		t.Fatal(err)
	}
	err = inputDef.Open()
	if err != nil {
		t.Fatal(err)
	}
}

// Verify the InputDefinition Encoding to the internal format
func TestInputDefinition_Encoding(t *testing.T) {
	inputBody := []byte(`
			{
			"frames":[{
				"name":"event-time",
				"options":{
					"timeQuantum": "YMD",
					"inverseEnabled": true,
					"cacheType": "ranked"
				}
			}],
			"fields": [
				{
					"name": "id",
					"primaryKey": true
				},
				{
					"name": "cabType",
					"actions": [
						{
							"frame": "cab-type",
							"valueDestination": "mapping",
							"valueMap": {
								"Green": 1,
								"Yellow": 2
							}
						}
					]
				}
			]
		}`)
	var def pilosa.InputDefinitionInfo
	err := json.Unmarshal(inputBody, &def)
	if err != nil {
		t.Fatal(err)
	}

	internalDef := def.Encode()

	if internalDef.Frames[0].Name != "event-time" {
		t.Fatalf("unexpected frame: %v", internalDef)
	} else if internalDef.Frames[0].Meta.CacheType != "ranked" {
		t.Fatalf("unexpected frame meta data: %v", internalDef)
	} else if len(internalDef.Fields) != 2 {
		t.Fatalf("unexpected number of Fields: %d", len(internalDef.Fields))
	} else if len(internalDef.Fields[1].Actions) != 1 {
		t.Fatalf("unexpected number of Actions: %v", internalDef.Fields[1].Actions)
	} else if internalDef.Fields[1].Actions[0].ValueDestination != "mapping" {
		t.Fatalf("unexpected ValueDestination: %v", internalDef.Fields[1].Actions[0])
	}
}

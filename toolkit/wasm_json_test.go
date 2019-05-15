package toolkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"testing"
)

var expectedJson = []JSON{
	{
		"Name":    "preramble",
		"Magic":   []byte{0, 97, 115, 109},
		"Version": []byte{1, 0, 0, 0},
	},
	{
		"Name":        "custom",
		"SectionName": "a custom section",
		"Payload":     "this is the payload",
	},
}

func TestCustomSection(t *testing.T) {
	wasm, err := ioutil.ReadFile(path.Join("test", "customSection.wasm"))
	assert.Nil(t, err)
	jsonObj := Wasm2Json(wasm)
	assert.Equal(t, true, assert.ObjectsAreEqual(expectedJson, jsonObj))

	generatedWasm := Json2Wasm(jsonObj)
	assert.Equal(t, 0, bytes.Compare(generatedWasm, wasm))
	assert.Equal(t, true, assert.ObjectsAreEqual(expectedJson, jsonObj))
}

//func TestBuild(t *testing.T) {
//	dirName := path.Join("test", "wasm")
//	dir, err := ioutil.ReadDir(dirName)
//	assert.Nil(t, err)
//
//	for _, fi := range dir {
//		if fi.IsDir() {
//			continue
//		}
//		wasm, err := ioutil.ReadFile(path.Join(dirName, fi.Name()))
//		assert.Nil(t, err)
//
//	}
//}
func readWasmModule(path string) ([]JSON, error) {
	var jsonArr []JSON
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(jsonData, &jsonArr); err != nil {
		return nil, err
	}
	return jsonArr, nil
}

func TestBasicTest(t *testing.T) {
	dirName := path.Join("test", "wasm")
	dir, err := ioutil.ReadDir(dirName)
	assert.Nil(t, err)
	failed := 0
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}

		wasm, err := ioutil.ReadFile(path.Join(dirName, fi.Name()))
		assert.Nil(t, err)

		jsonObj := Wasm2Json(wasm)
		wasmBin := Json2Wasm(jsonObj)

		if !assert.Equal(t, 0, bytes.Compare(wasm, wasmBin)) {
			failed += 1
			//fmt.Printf("%#v\n", jsonObj)
		}
	}

	fmt.Printf("total failed case %d\n", failed)

	dirName = path.Join("test", "json")
	dir, err = ioutil.ReadDir(dirName)
	assert.Nil(t, err)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}

		jsonObj, err := readWasmModule(path.Join(dirName, fi.Name()))
		assert.Nil(t, err)
		fmt.Printf("%#v\n", jsonObj)

		wasm := Json2Wasm(jsonObj)
		jsonObj2 := Wasm2Json(wasm)
		fmt.Printf("%#v\n", jsonObj2)

		assert.Equal(t, true, assert.ObjectsAreEqual(jsonObj, jsonObj2))
	}

}

func TestText2Json(t *testing.T) {
	text := "i32.const 32 drop"
	json := Text2Json(text)

	expected := []JSON{
		{
			"Name":       "const",
			"ReturnType": "i32",
			"Immediates": "32",
		}, {
			"Name": "drop",
		},
	}

	assert.Equal(t, true, assert.ObjectsAreEqual(expected, json))

	text = "br_table 0 0 0 0 i64.const 24"
	json = Text2Json(text)

	expected = []JSON{
		{
			"Name":       "br_table",
			"Immediates": []string{"0", "0", "0", "0"},
		}, {
			"ReturnType": "i64",
			"Name":       "const",
			"Immediates": "24",
		},
	}

	assert.Equal(t, true, assert.ObjectsAreEqual(expected, json))

	text = "call_indirect 1 i64.const 24"
	json = Text2Json(text)

	expected = []JSON{
		{
			"Name": "call_indirect",
			"Immediates": JSON{
				"Index":    "1",
				"Reserved": 0,
			},
		}, {
			"ReturnType": "i64",
			"Name":       "const",
			"Immediates": "24",
		},
	}

	assert.Equal(t, true, assert.ObjectsAreEqual(expected, json))

	text = "i32.load 0 1 i64.const 24"
	json = Text2Json(text)

	expected = []JSON{
		{
			"Name":       "load",
			"ReturnType": "i32",
			"Immediates": JSON{
				"Flags":  "0",
				"Offset": "1",
			},
		}, {
			"ReturnType": "i64",
			"Name":       "const",
			"Immediates": "24",
		},
	}

	assert.Equal(t, true, assert.ObjectsAreEqual(expected, json))
}

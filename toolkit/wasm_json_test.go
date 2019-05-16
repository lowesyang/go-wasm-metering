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
		"name":    "preramble",
		"magic":   []byte{0, 97, 115, 109},
		"version": []byte{1, 0, 0, 0},
	},
	{
		"name":        "custom",
		"section_name": "a custom section",
		"payload":     "this is the payload",
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

	//dirName = path.Join("test", "json")
	//dir, err = ioutil.ReadDir(dirName)
	//assert.Nil(t, err)
	//for _, fi := range dir {
	//	if fi.IsDir() {
	//		continue
	//	}
	//
	//	jsonObj, err := readWasmModule(path.Join(dirName, fi.Name()))
	//	assert.Nil(t, err)
	//	fmt.Printf("%#v\n", jsonObj)
	//
	//	wasm := Json2Wasm(jsonObj)
	//	jsonObj2 := Wasm2Json(wasm)
	//	fmt.Printf("%#v\n", jsonObj2)
	//
	//	assert.Equal(t, true, assert.ObjectsAreEqual(jsonObj, jsonObj2))
	//}

}

func TestText2Json(t *testing.T) {
	text := "i32.const 32 drop"
	json := Text2Json(text)

	expected := []JSON{
		{
			"name":        "const",
			"return_type": "i32",
			"immediates":  "32",
		}, {
			"name": "drop",
		},
	}

	assert.Equal(t, true, assert.ObjectsAreEqual(expected, json))

	text = "br_table 0 0 0 0 i64.const 24"
	json = Text2Json(text)

	expected = []JSON{
		{
			"name":       "br_table",
			"immediates": []string{"0", "0", "0", "0"},
		}, {
			"return_type": "i64",
			"name":       "const",
			"immediates": "24",
		},
	}

	assert.Equal(t, true, assert.ObjectsAreEqual(expected, json))

	text = "call_indirect 1 i64.const 24"
	json = Text2Json(text)

	expected = []JSON{
		{
			"name": "call_indirect",
			"immediates": JSON{
				"index":    "1",
				"reserved": 0,
			},
		}, {
			"return_type": "i64",
			"name":        "const",
			"immediates":  "24",
		},
	}

	assert.Equal(t, true, assert.ObjectsAreEqual(expected, json))

	text = "i32.load 0 1 i64.const 24"
	json = Text2Json(text)

	expected = []JSON{
		{
			"name":       "load",
			"return_type": "i32",
			"immediates": JSON{
				"flags":  "0",
				"offset": "1",
			},
		}, {
			"return_type": "i64",
			"name":       "const",
			"immediates": "24",
		},
	}

	assert.Equal(t, true, assert.ObjectsAreEqual(expected, json))
}

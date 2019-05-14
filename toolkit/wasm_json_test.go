package toolkit

import (
	"bytes"
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
}

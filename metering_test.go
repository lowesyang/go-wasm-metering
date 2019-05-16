package go_wasm_metering

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yyh1102/go-wasm-metering/toolkit"
	"io/ioutil"
	"path"
	"testing"
)

var (
	defaultCostTable toolkit.JSON
	defaultCTPath    = path.Join("test", "defaultCostTable.json")
)

func init() {
	defaultCostT, err := readCostTable(defaultCTPath)
	if err != nil {
		panic(err)
	}

	defaultCostTable = defaultCostT
}

func readCostTable(path string) (toolkit.JSON, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	jsonObj := make(toolkit.JSON)
	err = json.Unmarshal(data, &jsonObj)
	if err != nil {
		return nil, err
	}
	return jsonObj, nil
}

func readWasmModule(path string) ([]toolkit.JSON, error) {
	var jsonArr []toolkit.JSON
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(jsonData, &jsonArr); err != nil {
		return nil, err
	}
	return jsonArr, nil
}

func TestBasic(t *testing.T) {
	wasm, err := ioutil.ReadFile(path.Join("test", "in", "wasm", "basic.wasm"))
	assert.Nil(t, err)

	meteredWasm, err := MeterWASM(wasm, &Options{
		CostTable: defaultCTPath,
	})
	assert.Nil(t, err)
	meteredJson := toolkit.Wasm2Json(meteredWasm)

	expectedWasm, err := ioutil.ReadFile(path.Join("test", "expected-out", "wasm", "basic.wasm"))
	assert.Nil(t, err)
	expectedJson := toolkit.Wasm2Json(expectedWasm)
	//fmt.Printf("%#v\n%#v\n", meteredJson, expectedJson)
	assert.Equal(t, true, assert.ObjectsAreEqual(meteredJson, expectedJson))
	assert.Equal(t, 0, bytes.Compare(meteredWasm, expectedWasm))
	//entries1 := meteredJson[1]["entries"].([]toolkit.TypeEntry)
	//entries2 := meteredJson[2]["entries"].([]toolkit.TypeEntry)
	//assert.Equal(t, "metering", entries2[0]["module_str"])
	//assert.Equal(t, "usegas", entries2[0]["field_str"].(string))
	//assert.Equal(t, "i32", entries1[1]["params"].([]interface{})[0].(string))
}

func TestBasicMeteringTests(t *testing.T) {
	dirName := path.Join("test", "in")
	dir, err := ioutil.ReadDir(path.Join(dirName, "wasm"))
	assert.Nil(t, err)
	for _, file := range dir {
		// read wasm json.
		wasm, err := ioutil.ReadFile(path.Join(dirName, "wasm", file.Name()))
		assert.Nil(t, err)

		module := toolkit.Wasm2Json(wasm)

		// read cost table json.
		costTable, err := readCostTable(path.Join(dirName, "costTables", file.Name()))
		if err != nil {
			costTable = defaultCostTable
		}
		metering := Metering{
			costTable: costTable,
			opts: Options{
				ModuleStr: defaultModuleStr,
				FieldStr:  defaultFieldStr,
				MeterType: defaultMeterType,
			},
		}
		//fmt.Printf("%s %#v\n", file.Name(), module)
		meteredModule, err := metering.meterJSON(module)
		if err != nil {
			assert.Equal(t, "basic+import.wasm", file.Name())
			continue
		}
		//fmt.Printf("%s old %#v\n", file.Name(), meteredModule)

		expectedWasm, err := ioutil.ReadFile(path.Join("test", "expected-out", "wasm", file.Name()))
		assert.Nil(t, err)
		expectedJson := toolkit.Wasm2Json(expectedWasm)
		//fmt.Printf("%s exp %#v\n", file.Name(), expectedJson)

		if !assert.Equal(t, true, assert.ObjectsAreEqual(meteredModule, expectedJson)) {
			fmt.Printf("%#v\n%#v\n", meteredModule, expectedJson)
		}
	}

	//fmt.Printf("Basic metering tests failed cases %d", failed)
}

func TestWasm(t *testing.T) {
	module, err := readWasmModule(path.Join("test", "in", "json", "basic.wast.json"))
	assert.Nil(t, err)

	wasm := toolkit.Json2Wasm(module)
	meteredWasm, err := MeterWASM(wasm, &Options{
		MeterType: "i32",
		FieldStr:  "test",
		ModuleStr: "test",
	})
	assert.Nil(t, err)

	meteredJson := toolkit.Wasm2Json(meteredWasm)
	entries1 := meteredJson[1]["Entries"].([]interface{})
	entries2 := meteredJson[2]["Entries"].([]interface{})
	assert.Equal(t, "test", entries2[0].(toolkit.JSON)["ModuleStr"].(string))
	assert.Equal(t, "test", entries2[0].(toolkit.JSON)["FieldStr"].(string))
	assert.Equal(t, "i32", entries1[1].(toolkit.JSON)["Params"].([]interface{})[0].(string))
}

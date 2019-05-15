package go_wasm_metering

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yyh1102/go-wasm-metering/toolkit"
	"io/ioutil"
	"path"
	"testing"
)

var defaultCostTable toolkit.JSON

func init() {
	defaultCostT, err := readCostTable(path.Join("test", "defaultCostTable.json"))
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
	module, err := readWasmModule(path.Join("test", "in", "json", "basic.wast.json"))
	assert.Nil(t, err)
	wasm := toolkit.Json2Wasm(module)
	fmt.Printf("%#v\n", toolkit.Wasm2Json(wasm))

	meteredWasm, err := MeterWASM(wasm, nil)
	assert.Nil(t, err)
	meteredJson := toolkit.Wasm2Json(meteredWasm)

	fmt.Printf("%#v\n", meteredJson)
	entries1 := meteredJson[1]["Entries"].([]interface{})
	entries2 := meteredJson[2]["Entries"].([]interface{})
	assert.Equal(t, "metering", entries2[0].(toolkit.JSON)["ModuleStr"].(string))
	assert.Equal(t, "usegas", entries2[0].(toolkit.JSON)["FieldStr"].(string))
	assert.Equal(t, "i32", entries1[1].(toolkit.JSON)["Params"].([]interface{})[0].(string))
}

func TestBasicMeteringTests(t *testing.T) {
	dirName := path.Join("test", "in")
	dir, err := ioutil.ReadDir(dirName)
	assert.Nil(t, err)
	failed := 0
	for _, file := range dir {
		// read wasm json.
		module, err := readWasmModule(path.Join(dirName, "json", file.Name()))
		assert.Nil(t, err)

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
		meteredModule, err := metering.meterJSON(module)
		assert.Nil(t, err)

		expectedModule, err := readWasmModule(path.Join("test", "expected-out", "json", file.Name()))
		assert.Nil(t, err)
		if !assert.Equal(t, true, assert.ObjectsAreEqual(expectedModule, meteredModule)) {
			//assert.Equal(t, "basic+import.wasm.json", file.Name())
			failed += 1
		}
	}

	assert.Equal(t, 0, failed)
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

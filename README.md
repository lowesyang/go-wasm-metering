# go-wasm-metering
Gas metering injector for eWASM in Golang.

## Install

`go get -u github.com/yyh1102/go-wasm-metering`

## Usage

Inject meter func to wasm.

```Go
package main

import "github.com/yyh1102/go-wasm-metering"

func main(){
	wasm, err:=ioutil.ReadFile("xxx")
	if err!=nil{
		panic(err)
	}
	
	opts := &metering.Options{}
	
	meterWasm:=metering.MeterWasm(wasm,opts)
	fmt.Println(meterWasm)
}

```

## License
Apache-2.0

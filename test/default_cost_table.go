package test

import "github.com/yyh1102/go-wasm-metering/toolkit"

var DefaultCostTable = toolkit.JSON{
	"start": 1,
	"type": toolkit.JSON{
		"params": toolkit.JSON{
			"DEFAULT": 1,
		},
		"return_type": toolkit.JSON{
			"DEFAULT": 1,
		},
	},
	"import": 1,
	"code": toolkit.JSON{
		"locals": toolkit.JSON{
			"DEFAULT": 1,
		},
		"code": toolkit.JSON{
			"DEFAULT": 1,
		},
	},
	"data": 0,
}

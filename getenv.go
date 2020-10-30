package main

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/skeris/flat-grabber/app"
)

func getEnv() app.Options {
	var res app.Options

	r := reflect.ValueOf(res)
	for i := 0; i < r.NumField(); i++ {
		v := r.Type().Field(i).Tag.Get(app.ENV)
		d := r.Type().Field(i).Tag.Get(app.DEFAULT)

		env, ok := os.LookupEnv(v)
		var val string
		if ok {
			val = env
		} else {
			val = d
		}

		switch t := r.Type().Field(i).Type.Name(); t {
		case "string":
			reflect.ValueOf(&res).Elem().Field(i).SetString(val)
		case "bool":
			if strings.ToLower(val) == "true" {
				reflect.ValueOf(&res).Elem().Field(i).SetBool(true)
			} else {
				reflect.ValueOf(&res).Elem().Field(i).SetBool(false)
			}
		case "uint64":
			num, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				panic(err)
			}
			reflect.ValueOf(&res).Elem().Field(i).SetUint(num)
		default:
			panic(errors.New("Something strange happend: " + t))
		}
	}

	return res
}

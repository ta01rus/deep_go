package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

type Car struct {
	Mark         string            `properties:"mark"`
	Producer     string            `properties:",omitempty"`
	ProducerDate time.Time         `properties:""`
	Power        int               `properties:"power,omitempty"`
	Chars        map[string]string `properties:"chars,omitempty"`
	Owners       []string          `properties:"owners,omitempty"`
}

func Serialize(obj any) string {

	objValueType := reflect.ValueOf(obj)
	objType := reflect.TypeOf(obj)

	switch objValueType.Kind() {
	case reflect.Int64, reflect.Int, reflect.Int32, reflect.Int8:
		return strconv.Itoa(int(objValueType.Int()))
	case reflect.Float64, reflect.Float32:
		return strconv.FormatFloat(objValueType.Float(), 'f', 6, 64)
	case reflect.Bool:
		if objValueType.Bool() {
			return "true"
		}
		return "false"
	case reflect.Map:
		// тут могут проблемы с порядком в мапе и тестом
		parts := make([]string, 0, len(objValueType.MapKeys()))
		mit := objValueType.MapRange()
		for mit.Next() {
			parts = append(parts, mit.Key().String()+":"+mit.Value().String())
		}
		return strings.Join(parts, ",")
	case reflect.Slice:
		parts := make([]string, 0, objValueType.Len())
		for i := 0; i < objValueType.Len(); i++ {
			el := objValueType.Index(i)
			val := Serialize(el.Interface())
			parts = append(parts, val)
		}
		return strings.Join(parts, ",")
	case reflect.Pointer:
		return Serialize(objValueType.Elem().Interface())
	case reflect.Struct:
		switch t := objValueType.Interface().(type) {
		case time.Time:
			return t.Format("2006-01-02 15:04:05")
		default:
			parts := make([]string, 0, objType.NumField())
			for i := 0; i < objType.NumField(); i++ {
				fieldValue, fieldType := objValueType.Field(i), objType.Field(i)

				key := fieldType.Name
				value := Serialize(fieldValue.Interface())

				if tags, ok := fieldType.Tag.Lookup("properties"); ok {
					options := strings.Split(tags, ",")

					if len(options) > 0 && options[0] != "" {
						key = options[0]
					}
					switch {
					case value == "" && len(options) < 1:
						panic(fmt.Errorf("%s is empty ", fieldType.Name))
					case value == "" && len(options) > 1 && options[1] != "omitempty":
						panic(fmt.Errorf("%s is empty ", fieldType.Name))
					case value == "" && len(options) > 1 && options[1] == "omitempty":
						continue
					}

				}
				parts = append(parts, key+"="+value)
			}
			return strings.Join(parts, "\n")
		}

	case reflect.String:
		return objValueType.String()
	default:
		panic(fmt.Errorf("can't convert to string type"))
	}
}

// go test -v homework_test.go
func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestCarSerialization(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2000-04-01")
	tests := map[string]struct {
		car    Car
		result string
	}{

		"test case with empty name": {
			car: Car{
				Mark:         "Mark II",
				Producer:     "Toyota",
				ProducerDate: date,
			},
			result: "mark=Mark II\nProducer=Toyota\nProducerDate=2000-04-01 00:00:00\npower=0",
		},
		"test continue omitempty ": {
			car: Car{
				Mark:         "Mark II",
				ProducerDate: date,
			},
			result: "mark=Mark II\nProducerDate=2000-04-01 00:00:00\npower=0",
		},
		"test map Serializer": {
			car: Car{
				Mark:         "Mark II",
				ProducerDate: date,
				Chars: map[string]string{
					"color": "red",
					"speed": "320",
				},
			},
			result: "mark=Mark II\nProducerDate=2000-04-01 00:00:00\npower=0\nchars=color:red,speed:320",
		},
		"test map slice": {
			car: Car{
				Mark:         "Mark II",
				ProducerDate: date,
				Chars: map[string]string{
					"color": "red",
					"speed": "320",
				},
				Owners: []string{"Иванов", "Петров", "Сидоров"},
			},
			result: "mark=Mark II\nProducerDate=2000-04-01 00:00:00\npower=0\nchars=color:red,speed:320\nowners=Иванов,Петров,Сидоров",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.car)
			assert.Equal(t, test.result, result)
		})
	}
}

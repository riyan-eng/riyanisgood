package riyanisgood

import (
    "reflect"
	"strings"

	"encoding/json"
	"net/url"

	"github.com/thedevsaddam/govalidator"
)

type validatorInterface interface{
    ValidateStruct(dataStruct any) (validErrorSlice url.Values)
}

type validatorStruct struct{

}

func NewValidation() validatorInterface {
    return &validStruct{}
}

func (v *validStruct) ValidateStruct(dataStruct any) (validErrorSlice url.Values) {
	rv := reflect.ValueOf(dataStruct)
	rt := rv.Type()
	var validRulesSlice []validStruct
	var validMessagesSlice []validStruct
	for i := 0; i < rt.NumField(); i++ {
		if value, ok := rt.Field(i).Tag.Lookup("valid"); ok {
			validRulesSlice = append(validRulesSlice, validStruct{
				Key:   rt.Field(i).Tag.Get("json"),
				Value: strings.Split(value, ";"),
			})
		}
		if value, ok := rt.Field(i).Tag.Lookup("valid_message"); ok {
			validMessagesSlice = append(validMessagesSlice, validStruct{
				Key:   rt.Field(i).Tag.Get("json"),
				Value: strings.Split(value, ";"),
			})
		}
	}
	validRulesMap := convertStructToMap(validRulesSlice)
	validMessagesMap := convertStructToMap(validMessagesSlice)
	var validDataMap map[string]interface{}
	data, _ := json.Marshal(dataStruct)
	json.Unmarshal(data, &validDataMap)
	validErrorSlice = generateErrorSlice(validRulesMap, validMessagesMap, validDataMap)
	return
}

type validStruct struct {
	Key   string
	Value []string
}

func convertStructToMap(validSlice []validStruct) (validMap map[string][]string) {
	validMap = make(map[string][]string)
	for _, val := range validSlice {
		validMap[val.Key] = append(validMap[val.Key], val.Value...)
	}
	return
}

func generateErrorSlice(rules map[string][]string, messages map[string][]string, data map[string]interface{}) (validErrorSlice url.Values) {
	opts := govalidator.Options{
		Data:     &data,
		Rules:    rules,
		Messages: messages,
	}
	v := govalidator.New(opts)
	validErrorSlice = v.ValidateStruct()
	return
}

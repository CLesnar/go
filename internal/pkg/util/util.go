package util

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	cacheStructFieldTagMap = map[string]string{}
)

func lookupStructFieldTag(model interface{}, field interface{}, tagsKey string) (string, bool) {
	if mv, fv := reflect.ValueOf(model), reflect.ValueOf(field); fv.Kind() != reflect.Pointer || mv.Kind() != reflect.Pointer {
		return "", false
	} else {
		fieldOffset := fv.Pointer() - mv.Pointer()
		s := fmt.Sprintf("%t.%v.%s", model, fieldOffset, tagsKey)
		tag, ok := cacheStructFieldTagMap[s]
		return tag, ok
	}
}

func cacheStructFieldTag(model interface{}, field interface{}, tagsKey string, tag string) {
	if mv, fv := reflect.ValueOf(model), reflect.ValueOf(field); fv.Kind() != reflect.Pointer || mv.Kind() != reflect.Pointer {
		return
	} else {
		fieldOffset := fv.Pointer() - mv.Pointer()
		s := fmt.Sprintf("%t.%v.%s", model, fieldOffset, tagsKey)
		cacheStructFieldTagMap[s] = tag
	}
}

func FindFieldTag(model interface{}, field interface{}, tagsKey string) (string, error) {
	if tag, ok := lookupStructFieldTag(model, field, tagsKey); ok {
		return tag, nil
	}
	mv := reflect.ValueOf(model).Elem()
	fv := reflect.ValueOf(field)
	len := mv.NumField()
	for i := 0; i < len; i++ {
		f := mv.Field(i).Addr()
		if f.Pointer() == fv.Pointer() {
			if f.Type() == fv.Type() {
				if tagAndOptions, ok := mv.Type().Field(i).Tag.Lookup(tagsKey); ok {
					tag := strings.Split(tagAndOptions, ",")[0]
					cacheStructFieldTag(model, field, tagsKey, tag)
					return tag, nil
				}
			}
		}
	}
	return "", fmt.Errorf("failed to get tag for field '%T:%+v' from struct '%T:%+v'", field, field, model, model)
}

func Reference[T any](v T) *T {
	return &v
}

func Contains[T any](list []T, item T, comparator func(a T, b T) int) bool {
	var e T
	for _, e = range list {
		if comparator(item, e) == 0 {
			return true
		}
	}
	return false
}

func GetTagsByOption(model interface{}, tagsReturnKey string, tagsFindKey string, tagsFindOption string, tagsSkipKey string, skipTagOrOption string) []string {
	rv := reflect.TypeOf(model).Elem()
	tags := []string{}
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		if len(tagsSkipKey) > 0 && len(skipTagOrOption) > 0 && strings.Contains(f.Tag.Get(tagsSkipKey), skipTagOrOption) {
			continue
		}
		returnTagAndOptions := f.Tag.Get(tagsReturnKey)
		returnTag := strings.Split(returnTagAndOptions, ",")[0]
		if len(tagsFindKey) < 1 || len(tagsFindOption) < 1 || strings.Contains(f.Tag.Get(tagsFindKey), tagsFindOption) {
			tags = append(tags, returnTag)
		}
	}
	return tags
}

func GetTags(model interface{}, tagsKey string) []string {
	return GetTagsByOption(model, tagsKey, "", "", "", "")
}

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

func ConstantRef[T any](v T) *T {
	return &v
}

func PrepareFilterText(s string) string {
	str := strings.ReplaceAll(s, "'", "''")
	return str
}

func StringifyUnknown(value interface{}, typeValueFormatOverrideFunc func(interface{}) string) string {
	if value == nil {
		if typeValueFormatOverrideFunc != nil {
			return typeValueFormatOverrideFunc(value)
		}
		return fmt.Sprintf("%v", value)
	}
	v := reflect.ValueOf(value)
	if v.Type().Kind() == reflect.Slice || v.Type().Kind() == reflect.Array {
		length, values := v.Len(), ""
		for i := 0; i < length; i++ {
			vi := v.Index(i)
			if vi.Kind() == reflect.Interface || vi.Kind() == reflect.Pointer {
				vi = vi.Elem()
			}
			isString := vi.Kind() == reflect.String
			elemI := v.Index(i).Interface()
			elem := fmt.Sprintf("%v", elemI)
			if elemI == nil {
				continue
			} else if typeValueFormatOverrideFunc != nil {
				elem = typeValueFormatOverrideFunc(elemI)
			} else if isString {
				elem = "'" + PrepareFilterText(elem) + "'"
			}
			values = values + "," + elem
		}
		return values[1:]
	} else {
		var elem string
		if typeValueFormatOverrideFunc != nil {
			elem = typeValueFormatOverrideFunc(value)
		} else {
			elem = fmt.Sprintf("%v", value)
			if v.Kind() == reflect.String {
				elem = "'" + PrepareFilterText(elem) + "'"
			}
		}
		return elem
	}
}

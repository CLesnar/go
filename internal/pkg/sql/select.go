package sql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/CLesnar/go/internal/pkg/util"
)

type Params interface {
	BuildQuery() *SelectInfo
}

type Page[T any] struct {
	Data             []T
	Limit            int32
	Offset           int32
	TotalCount       int32
	QueryExecutionId *string
	NextToken        *string
}

type SelectInfo struct {
	CriteriaModel     Params
	TableModel        interface{}
	Table             string
	DistinctColumn    string
	Columns           []string
	WhereClause       []string
	Err               error
	Query             string
	Joins             []string
	CustomQuery       *string
	Limit             int32
	Offset            int32
	DefaultSortColumn string
	SortColumn        string
	SortDirection     string
	QueryExecutionId  *string
	NextToken         *string
}

type ColumnFilter struct {
	Column    string      `json:"column"`
	Operation string      `json:"operation"`
	Value     interface{} `json:"value"`
}

type FieldOp func(f interface{}) string

type FilterFormat struct {
	Format          string
	FormatValueFunc func(interface{}) string
}

type AddCriteriaConfig struct {
	TableAlias                      *string
	FieldTagsExclude                []string
	fieldOverrideColumnOperatorsMap map[uintptr]string
}

type ColumnFilterConfig struct {
	TableAlias *string
}

var (
	tableMap = map[string]string{
		// fmt.Sprintf("%T", <model>): "<table>"
	}

	criteriaFields = []string{"selectcolumns", "limit", "offset", "sort", "queryExecutionId", "nextToken", "columnFiltersJson"}

	columnFilterOperationMap = map[string]FilterFormat{
		"gt":    {Format: "%s > %v"},
		"gte":   {Format: "%s >= %v"},
		"lt":    {Format: "%s < %v"},
		"lte":   {Format: "%s <= %v"},
		"e":     {Format: "%s = %v"},
		"ne":    {Format: "%s != %v"},
		"in":    {Format: "%s in (%v)"},
		"nin":   {Format: "%s not in (%v)"},
		"like":  {Format: "regexp_like(%s, %v)"},
		"ilike": {Format: "regexp_like(lower(%s), %v)", FormatValueFunc: func(val interface{}) string { return strings.ToLower(fmt.Sprintf("%v", val)) }},
	}

	criteriaFieldTypeOperationMap = map[string]string{
		"string":  "ilike",
		"default": "e",
	}

	timestampAwsFormatFunc = func(v interface{}) string { return fmt.Sprintf("timestamp '%v'", v) }
	datestampAwsFormatFunc = func(v interface{}) string { return fmt.Sprintf("date '%v'", v) }
	columnTypeMap          = map[string]FilterFormat{
		fmt.Sprintf("%T", time.Time{}): {Format: "timestamp '%v'", FormatValueFunc: timestampAwsFormatFunc},
	}

	TableMap = map[string]string{}
)

func (cf *ColumnFilter) Validate(model interface{}) (string, string, error) {
	if cf == nil {
		return "", "", errors.New("column filter is nil")
	}
	modelColumns := util.GetTags(model, "json")
	if !util.Contains(modelColumns, cf.Column, strings.Compare) {
		return "", "", fmt.Errorf("column '%s' is not found in model '%T'", cf.Column, model)
	}
	if _, ok := columnFilterOperationMap[cf.Operation]; !ok {
		return "", "", errors.New("operation not found")
	}
	table, ok := TableMap[fmt.Sprintf("%T", model)]
	if !ok {
		return "", "", fmt.Errorf("failed to find table name for type '%T'", model)
	}
	joinTable, modelFieldType := "", "" // TODO util.GetFieldAndJoinTableByTag(model, "json", cf.Column)
	if len(joinTable) > 0 {
		table = joinTable
	}
	modelFieldType = strings.ReplaceAll(modelFieldType, "*", "")
	if modelFieldType == "string" {
		if reflect.ValueOf(cf.Value).Type().String() != "string" {
			cf.Value = fmt.Sprintf("%v", cf.Value)
		}
	} else {
		if stringOps := []string{"like", "ilike"}; util.Contains(stringOps, cf.Operation, strings.Compare) {
			return "", "", fmt.Errorf("column '%s', type '%s', does not allow operations '%s'", cf.Column, modelFieldType, strings.Join(stringOps, ","))
		}
	}
	return table, modelFieldType, nil
}

func (cf ColumnFilter) SliceExpected() bool {
	return util.Contains([]string{"in", "nin"}, cf.Operation, strings.Compare)
}

func (cf ColumnFilter) GetIsNilAllowedAndFormat() (string, error) {
	if cf.Value != nil {
		return "", nil
	}
	switch cf.Operation {
	case "e":
		return "%s is null", nil
	case "ne":
		return "%s is not null", nil
	}
	return "", fmt.Errorf("null not allowed in '%s' operation", cf.Operation)
}

func GetFormatedValue(cf ColumnFilter, fieldType string, formatColumnFilter FilterFormat) string {
	var value string // Get value formated
	if typeFilter, ok := columnTypeMap[fieldType]; ok && !cf.SliceExpected() {
		value = fmt.Sprintf(typeFilter.Format, cf.Value)
		if typeFilter.FormatValueFunc != nil {
			value = typeFilter.FormatValueFunc(cf.Value)
		}
	} else {
		value = util.StringifyUnknown(cf.Value, typeFilter.FormatValueFunc)
		if formatColumnFilter.FormatValueFunc != nil {
			value = formatColumnFilter.FormatValueFunc(value)
		}
	}
	return value
}

func (s *SelectInfo) AddWhereClauseCustom(whereClauseFmt string, args ...interface{}) {
	if len(whereClauseFmt) > 0 {
		whereClause := fmt.Sprintf(whereClauseFmt, args...)
		s.WhereClause = append(s.WhereClause, whereClause)
	}
}

func (s *SelectInfo) AddWhereClauseFromColumnFilter(cf ColumnFilter, config *ColumnFilterConfig) {
	table, fieldType, err := cf.Validate(s.TableModel)
	if err != nil {
		s.Err = errors.Join(s.Err, err)
		return
	}
	column := fmt.Sprintf("%s.%s", table, cf.Column)
	if config != nil && config.TableAlias != nil {
		if len(*config.TableAlias) < 1 {
			column = cf.Column
		} else {
			column = fmt.Sprintf("%s.%s", *config.TableAlias, cf.Column)
		}
	}
	if nilFormat, err := cf.GetIsNilAllowedAndFormat(); err != nil {
		s.Err = errors.Join(s.Err, err)
		return
	} else if len(nilFormat) < 1 {
		formatColumnFilter := columnFilterOperationMap[cf.Operation]
		value := GetFormatedValue(cf, fieldType, formatColumnFilter)
		s.AddWhereClauseCustom(formatColumnFilter.Format, column, value)
	} else {
		s.AddWhereClauseCustom(nilFormat, column)
	}
}

func (s *SelectInfo) AddColumnFilters(columnFilters *[]ColumnFilter, config ...ColumnFilterConfig) {
	if columnFilters == nil {
		return
	}
	var cfg *ColumnFilterConfig
	if len(config) > 0 {
		cfg = &config[0]
	}
	for _, cf := range *columnFilters {
		s.AddWhereClauseFromColumnFilter(cf, cfg)
	}
}

func getTypeOperation(typeStr string) string {
	if operation, ok := criteriaFieldTypeOperationMap[typeStr]; ok {
		return operation
	} else {
		return criteriaFieldTypeOperationMap["default"]
	}
}

func (s *SelectInfo) AddCriteriaFieldOperatorOverride(cfg *AddCriteriaConfig, operation string, fields ...interface{}) {
	if cfg == nil {
		s.Err = errors.Join(s.Err, errors.New("add criteria config cannot be nil"))
		return
	}
	mv := reflect.ValueOf(s.CriteriaModel)
	if mv.Kind() != reflect.Pointer {
		s.Err = errors.Join(s.Err, errors.New("selectioninfo criteria model must be pointer"))
		return
	}
	if cfg.fieldOverrideColumnOperatorsMap == nil {
		cfg.fieldOverrideColumnOperatorsMap = map[uintptr]string{}
	}
	for _, field := range fields {
		fv := reflect.ValueOf(field)
		if fv.Kind() != reflect.Pointer {
			continue // err?
		}
		fieldOffset := fv.Pointer() - mv.Pointer()
		cfg.fieldOverrideColumnOperatorsMap[fieldOffset] = operation
	}
}

func (s *SelectInfo) AddCriteria(config ...AddCriteriaConfig) {
	cfg := AddCriteriaConfig{}
	if len(config) > 0 {
		cfg = config[0]
	}
	cmv := reflect.ValueOf(s.CriteriaModel).Elem()
	cmvp := reflect.ValueOf(s.CriteriaModel).Pointer()
	cfs := []ColumnFilter{}
	numFields := cmv.NumField()
	for i := 0; i < numFields; i++ {
		f := cmv.Field(i)
		if f.Type().Kind() == reflect.Pointer && f.IsNil() {
			continue
		}
		fTag, fTagAndOptions := "", cmv.Type().Field(i).Tag.Get("json")
		if len(fTagAndOptions) > 0 {
			fTag = strings.Split(fTagAndOptions, ",")[0]
		}
		operation, useOverrideOp, canOverride := "", false, f.CanAddr()
		var fieldOffset uintptr
		if canOverride {
			fieldOffset = f.Addr().Pointer() - cmvp
		}
		if operation, useOverrideOp = cfg.fieldOverrideColumnOperatorsMap[fieldOffset]; useOverrideOp && canOverride {
			// overriding operation for field
		} else {
			if util.Contains(criteriaFields, fTag, strings.Compare) {
				continue
			}
			if util.Contains(cfg.FieldTagsExclude, fTag, strings.Compare) {
				continue
			}
		}
		if f.Type().Kind() == reflect.Pointer {
			f = f.Elem()
		}
		if f.CanInterface() {
			if !useOverrideOp {
				operation = getTypeOperation(f.Type().String())
			}
			// can add field specific operator overrides here
			cf := ColumnFilter{Column: fTag, Operation: operation, Value: f.Interface()}
			cfs = append(cfs, cf)
		}
	}
	s.AddColumnFilters(&cfs, ColumnFilterConfig{TableAlias: cfg.TableAlias})
}

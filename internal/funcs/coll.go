package funcs

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hairyhenderson/gomplate/v4/conv"
	"github.com/hairyhenderson/gomplate/v4/internal/deprecated"
	"github.com/hairyhenderson/gomplate/v4/internal/texttemplate"

	"github.com/hairyhenderson/gomplate/v4/coll"
)

// CreateCollFuncs -
func CreateCollFuncs(ctx context.Context) map[string]any {
	f := map[string]any{}

	ns := &CollFuncs{ctx}
	f["coll"] = func() any { return ns }

	f["has"] = ns.Has
	f["slice"] = ns.deprecatedSlice
	f["dict"] = ns.Dict
	f["keys"] = ns.Keys
	f["values"] = ns.Values
	f["append"] = ns.Append
	f["prepend"] = ns.Prepend
	f["uniq"] = ns.Uniq
	f["reverse"] = ns.Reverse
	f["merge"] = ns.Merge
	f["sort"] = ns.Sort
	f["jsonpath"] = ns.JSONPath
	f["jq"] = ns.JQ
	f["flatten"] = ns.Flatten
	f["set"] = ns.Set
	f["unset"] = ns.Unset
	return f
}

// CollFuncs -
type CollFuncs struct {
	ctx context.Context
}

// Slice -
func (CollFuncs) Slice(args ...any) []any {
	return coll.Slice(args...)
}

// deprecatedSlice -
// Deprecated: use coll.Slice instead
func (f *CollFuncs) deprecatedSlice(args ...any) []any {
	deprecated.WarnDeprecated(f.ctx, "the 'slice' alias for coll.Slice is deprecated - use coll.Slice instead")
	return coll.Slice(args...)
}

// GoSlice - same as text/template's 'slice' function
func (CollFuncs) GoSlice(item reflect.Value, indexes ...reflect.Value) (reflect.Value, error) {
	return texttemplate.GoSlice(item, indexes...)
}

// Has -
func (CollFuncs) Has(in any, key string) bool {
	return coll.Has(in, key)
}

// Index returns the result of indexing the last argument with the preceding
// index keys. This is similar to the `index` built-in template function, but
// the arguments are ordered differently for pipeline compatibility. Also, this
// function is more strict, and will return an error when the value doesn't
// contain the given key.
func (CollFuncs) Index(args ...any) (any, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("wrong number of args: wanted at least 2, got %d", len(args))
	}

	item := args[len(args)-1]
	indexes := args[:len(args)-1]

	return coll.Index(item, indexes...)
}

// Dict -
func (CollFuncs) Dict(in ...any) (map[string]any, error) {
	return coll.Dict(in...)
}

// Keys -
func (CollFuncs) Keys(in ...map[string]any) ([]string, error) {
	return coll.Keys(in...)
}

// Values -
func (CollFuncs) Values(in ...map[string]any) ([]any, error) {
	return coll.Values(in...)
}

// Append -
func (CollFuncs) Append(v any, list any) ([]any, error) {
	return coll.Append(v, list)
}

// Prepend -
func (CollFuncs) Prepend(v any, list any) ([]any, error) {
	return coll.Prepend(v, list)
}

// Uniq -
func (CollFuncs) Uniq(in any) ([]any, error) {
	return coll.Uniq(in)
}

// Reverse -
func (CollFuncs) Reverse(in any) ([]any, error) {
	return coll.Reverse(in)
}

// Merge -
func (CollFuncs) Merge(dst map[string]any, src ...map[string]any) (map[string]any, error) {
	return coll.Merge(dst, src...)
}

// Sort -
func (CollFuncs) Sort(args ...any) ([]any, error) {
	var (
		key  string
		list any
	)
	if len(args) == 0 || len(args) > 2 {
		return nil, fmt.Errorf("wrong number of args: wanted 1 or 2, got %d", len(args))
	}
	if len(args) == 1 {
		list = args[0]
	}
	if len(args) == 2 {
		key = conv.ToString(args[0])
		list = args[1]
	}
	return coll.Sort(key, list)
}

// JSONPath -
func (CollFuncs) JSONPath(p string, in any) (any, error) {
	return coll.JSONPath(p, in)
}

// JQ -
func (f *CollFuncs) JQ(jqExpr string, in any) (any, error) {
	return coll.JQ(f.ctx, jqExpr, in)
}

// Flatten -
func (CollFuncs) Flatten(args ...any) ([]any, error) {
	if len(args) == 0 || len(args) > 2 {
		return nil, fmt.Errorf("wrong number of args: wanted 1 or 2, got %d", len(args))
	}

	list := args[0]

	var err error
	depth := -1
	if len(args) == 2 {
		depth, err = conv.ToInt(list)
		if err != nil {
			return nil, fmt.Errorf("wrong depth type: must be int, got %T (%+v)", list, list)
		}

		list = args[1]
	}

	return coll.Flatten(list, depth)
}

func pickOmitArgs(args ...any) (map[string]any, []string, error) {
	if len(args) <= 1 {
		return nil, nil, fmt.Errorf("wrong number of args: wanted 2 or more, got %d", len(args))
	}

	m, ok := args[len(args)-1].(map[string]any)
	if !ok {
		return nil, nil, fmt.Errorf("wrong map type: must be map[string]any, got %T", args[len(args)-1])
	}

	// special-case - if there's only one key and it's a slice, expand it
	if len(args) == 2 {
		if reflect.TypeOf(args[0]).Kind() == reflect.Slice {
			sl := reflect.ValueOf(args[0])
			expandedArgs := make([]any, sl.Len()+1)
			for i := range sl.Len() {
				expandedArgs[i] = sl.Index(i).Interface()
			}
			expandedArgs[len(expandedArgs)-1] = m
			args = expandedArgs
		}
	}

	keys := make([]string, len(args)-1)
	for i, v := range args[0 : len(args)-1] {
		k, ok := v.(string)
		if !ok {
			return nil, nil, fmt.Errorf("wrong key type: must be string, got %T (%+v)", args[i], args[i])
		}
		keys[i] = k
	}
	return m, keys, nil
}

// Pick -
func (CollFuncs) Pick(args ...any) (map[string]any, error) {
	m, keys, err := pickOmitArgs(args...)
	if err != nil {
		return nil, err
	}
	return coll.Pick(m, keys...), nil
}

// Omit -
func (CollFuncs) Omit(args ...any) (map[string]any, error) {
	m, keys, err := pickOmitArgs(args...)
	if err != nil {
		return nil, err
	}
	return coll.Omit(m, keys...), nil
}

// Set -
func (CollFuncs) Set(key string, value any, m map[string]any) (map[string]any, error) {
	m[key] = value

	return m, nil
}

// Unset -
func (CollFuncs) Unset(key string, m map[string]any) (map[string]any, error) {
	delete(m, key)

	return m, nil
}

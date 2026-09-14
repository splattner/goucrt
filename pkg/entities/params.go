package entities

// CommandParams is the params map an entity command handler receives, decoded from the remote's
// entity_command request. Convert with entities.CommandParams(params) inside a handler to use
// these typed accessors instead of an unchecked type assertion on the raw map - handy because
// encoding/json always decodes a JSON number into a float64 when the target is interface{}, never
// int or uint, which is a real source of "compiles fine, panics on the first real command from the
// remote" bugs when a handler asserts the wrong numeric type.
type CommandParams map[string]interface{}

// String returns params[key] as a string, and whether it was present as one.
func (p CommandParams) String(key string) (string, bool) {
	v, ok := p[key].(string)
	return v, ok
}

// Float64 returns params[key] as a float64, and whether it was present as a JSON number. This is
// how every JSON number decodes - prefer this (or Int) over asserting int/uint directly.
func (p CommandParams) Float64(key string) (float64, bool) {
	v, ok := p[key].(float64)
	return v, ok
}

// Int returns params[key] truncated to an int, and whether it was present as a JSON number.
func (p CommandParams) Int(key string) (int, bool) {
	v, ok := p[key].(float64)
	if !ok {
		return 0, false
	}
	return int(v), true
}

// Bool returns params[key] as a bool, and whether it was present as one.
func (p CommandParams) Bool(key string) (bool, bool) {
	v, ok := p[key].(bool)
	return v, ok
}

// Has reports whether key is present in params at all.
func (p CommandParams) Has(key string) bool {
	_, ok := p[key]
	return ok
}

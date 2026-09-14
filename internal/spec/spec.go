// Package spec parses the vendored Core-API AsyncAPI definition (see spec/core-api) so goucrt's own
// tests can check its hand-maintained protocol constants against the spec instead of relying on
// someone noticing drift by eye. See spec/core-api/README.md for how the vendored copy is updated.
//
// Only `features` and `device_class` are exposed here: they are the only entity properties the spec
// defines as a structured JSON Schema enum. Attribute names, state values and command identifiers are
// not present as schema in the YAML - they only exist in doc/entities/*.md's prose tables and worked
// examples, which pkg/entities/testdata's golden fixtures are extracted from instead.
package spec

import (
	_ "embed"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed core-api/integration-api/UCR-integration-asyncapi.yaml
var asyncAPIYAML []byte

// EntitySchema holds the parts of an entity's AsyncAPI schema goucrt can check itself against.
type EntitySchema struct {
	// Features is the entity's `features` enum, in spec order. Empty if the entity type has no
	// features (e.g. sensor).
	Features []string
	// DeviceClass is the entity's `device_class` enum, in spec order. Empty if the entity type has
	// no device_class property.
	DeviceClass []string
}

// EntitySchema returns the features and device_class enums for the given entity_type discriminator
// (e.g. "light", "media_player") as defined in the vendored AsyncAPI spec.
func EntityTypeSchema(entityType string) (EntitySchema, error) {
	doc, err := document()
	if err != nil {
		return EntitySchema{}, err
	}

	schemas, err := mapAt(doc, "components", "schemas")
	if err != nil {
		return EntitySchema{}, err
	}

	raw, ok := schemas[entityType]
	if !ok {
		return EntitySchema{}, fmt.Errorf("spec: no schema for entity_type %q", entityType)
	}
	entitySchema, ok := raw.(map[string]interface{})
	if !ok {
		return EntitySchema{}, fmt.Errorf("spec: schema for entity_type %q is not an object", entityType)
	}

	// Entity schemas are `allOf: [{$ref: .../entity}, {type: object, properties: {...}}]`; the
	// entity-specific properties (features, device_class, options) live in the second element.
	allOf, ok := entitySchema["allOf"].([]interface{})
	if !ok || len(allOf) < 2 {
		return EntitySchema{}, fmt.Errorf("spec: schema for entity_type %q has no allOf[1] extension object", entityType)
	}
	ext, ok := allOf[1].(map[string]interface{})
	if !ok {
		return EntitySchema{}, fmt.Errorf("spec: schema for entity_type %q: allOf[1] is not an object", entityType)
	}
	props, _ := ext["properties"].(map[string]interface{})

	return EntitySchema{
		Features:    enumAt(schemas, props, "features", "items", "enum"),
		DeviceClass: enumAt(schemas, props, "device_class", "enum"),
	}, nil
}

var parsed map[string]interface{}
var parseErr error
var parsedOnce bool

// document lazily parses the embedded spec YAML once per process.
func document() (map[string]interface{}, error) {
	if parsedOnce {
		return parsed, parseErr
	}
	parsedOnce = true

	var doc map[string]interface{}
	if err := yaml.Unmarshal(asyncAPIYAML, &doc); err != nil {
		parseErr = fmt.Errorf("spec: parse embedded AsyncAPI YAML: %w", err)
		return nil, parseErr
	}
	parsed = doc
	return parsed, nil
}

// mapAt walks a chain of nested map keys, returning an error naming the first key that either
// doesn't exist or isn't itself a map.
func mapAt(doc map[string]interface{}, path ...string) (map[string]interface{}, error) {
	cur := doc
	walked := ""
	for _, key := range path {
		raw, ok := cur[key]
		if !ok {
			return nil, fmt.Errorf("spec: missing %s%s", walked, key)
		}
		next, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("spec: %s%s is not an object", walked, key)
		}
		cur = next
		walked += key + "."
	}
	return cur, nil
}

// enumAt walks props[path[0]][path[1]]...["enum"] and returns the enum's string values, or nil if
// any step along the way is absent (a property with no enum, e.g. an entity with no device_class).
// Resolves a `$ref: '#/components/schemas/X'` into schemas[X] before indexing into it at each step -
// e.g. voice_assistant's `features` property is itself a $ref (to VoiceAssistantFeatures, an array
// of a further $ref to the VoiceAssistantFeature enum) rather than the inline
// `items: {enum: [...]}` every other entity type uses.
func enumAt(schemas map[string]interface{}, props map[string]interface{}, path ...string) []string {
	if props == nil {
		return nil
	}
	var cur interface{} = props
	for _, key := range path {
		m, ok := resolveRef(schemas, cur).(map[string]interface{})
		if !ok {
			return nil
		}
		cur, ok = m[key]
		if !ok {
			return nil
		}
	}
	rawList, ok := cur.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(rawList))
	for _, v := range rawList {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// resolveRef returns schemas[X] if node is a map whose only relevant content is a
// `$ref: '#/components/schemas/X'`, following further refs up to a small hop limit (in case the
// referenced schema is itself just another ref) in case of a cycle in the vendored spec. Any other
// node - including a map with a $ref pointing outside components/schemas, which the vendored spec
// never does - is returned unchanged.
func resolveRef(schemas map[string]interface{}, node interface{}) interface{} {
	const prefix = "#/components/schemas/"
	for range 10 {
		m, ok := node.(map[string]interface{})
		if !ok {
			return node
		}
		ref, ok := m["$ref"].(string)
		if !ok {
			return node
		}
		name, ok := strings.CutPrefix(ref, prefix)
		if !ok {
			return node
		}
		next, ok := schemas[name]
		if !ok {
			return node
		}
		node = next
	}
	return node
}

package entities

import "testing"

// CommandParams' accessors exist because encoding/json always decodes a JSON number into float64
// when unmarshaled into interface{} - never int or uint - so these tests focus on that boundary
// (present-as-wrong-type and absent both returning ok=false) rather than the happy path alone.

func TestCommandParamsFloat64(t *testing.T) {
	p := CommandParams{
		"hue":    float64(180),
		"source": "HDMI1",
	}

	if v, ok := p.Float64("hue"); !ok || v != 180 {
		t.Errorf("Float64(%q) = %v, %v, want 180, true", "hue", v, ok)
	}
	if _, ok := p.Float64("source"); ok {
		t.Errorf("Float64(%q) on a string value: ok = true, want false", "source")
	}
	if _, ok := p.Float64("missing"); ok {
		t.Errorf("Float64(%q) on an absent key: ok = true, want false", "missing")
	}
}

func TestCommandParamsInt(t *testing.T) {
	p := CommandParams{"brightness": float64(76.9)}

	if v, ok := p.Int("brightness"); !ok || v != 76 {
		t.Errorf("Int(%q) = %v, %v, want 76, true (truncated, not rounded)", "brightness", v, ok)
	}
	if _, ok := p.Int("missing"); ok {
		t.Errorf("Int(%q) on an absent key: ok = true, want false", "missing")
	}
}

func TestCommandParamsString(t *testing.T) {
	p := CommandParams{
		"source": "HDMI1",
		"volume": float64(50),
	}

	if v, ok := p.String("source"); !ok || v != "HDMI1" {
		t.Errorf("String(%q) = %v, %v, want HDMI1, true", "source", v, ok)
	}
	if _, ok := p.String("volume"); ok {
		t.Errorf("String(%q) on a float64 value: ok = true, want false", "volume")
	}
}

func TestCommandParamsBool(t *testing.T) {
	p := CommandParams{"muted": true}

	if v, ok := p.Bool("muted"); !ok || !v {
		t.Errorf("Bool(%q) = %v, %v, want true, true", "muted", v, ok)
	}
	if _, ok := p.Bool("missing"); ok {
		t.Errorf("Bool(%q) on an absent key: ok = true, want false", "missing")
	}
}

func TestCommandParamsHas(t *testing.T) {
	p := CommandParams{"hue": float64(0)}

	if !p.Has("hue") {
		t.Error("Has(\"hue\") = false, want true even though the value is the zero value")
	}
	if p.Has("missing") {
		t.Error("Has(\"missing\") = true, want false")
	}
}

package cal

import (
	"time"

	"testing"
)

func TestLogicAnd(t *testing.T) {
	tests := []struct {
		name   string
		inputs []any
		want   bool
	}{
		{"all true", []any{1, "a", true}, true},
		{"contains false", []any{1, 0, true}, false},
		{"contains nil", []any{1, nil, true}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LogicAnd(tt.inputs[0], tt.inputs[1:]...)
			if got != tt.want {
				t.Errorf("LogicAnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogicOr(t *testing.T) {
	tests := []struct {
		name   string
		inputs []any
		want   bool
	}{
		{"all false", []any{0, "", nil}, false},
		{"contains true", []any{0, "a", nil}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LogicOr(tt.inputs[0], tt.inputs[1:]...)
			if got != tt.want {
				t.Errorf("LogicOr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogicEq(t *testing.T) {
	tests := []struct {
		name    string
		a, b    any
		want    bool
		wantErr bool
	}{
		{"int equal", 1, 1, true, false},
		{"int not equal", 1, 2, false, false},
		{"string equal", "a", "a", true, false},
		{"string not equal", "a", "b", false, false},
		{"time equal", time.Time{}, time.Time{}, true, false},
		{"time not equal", time.Now(), time.Now().Add(time.Microsecond), false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LogicEq(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogicEq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LogicEq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogicNeq(t *testing.T) {
	tests := []struct {
		name    string
		a, b    any
		want    bool
		wantErr bool
	}{
		{"int equal", 1, 1, false, false},
		{"int not equal", 1, 2, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LogicNeq(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogicNeq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LogicNeq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogicGt(t *testing.T) {
	tests := []struct {
		name    string
		a, b    any
		want    bool
		wantErr bool
	}{
		{"int greater", 2, 1, true, false},
		{"int equal", 1, 1, false, false},
		{"int less", 1, 2, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LogicGt(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogicGt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LogicGt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogicGte(t *testing.T) {
	tests := []struct {
		name    string
		a, b    any
		want    bool
		wantErr bool
	}{
		{"int greater", 2, 1, true, false},
		{"int equal", 1, 1, true, false},
		{"int less", 1, 2, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LogicGte(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogicGte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LogicGte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogicLt(t *testing.T) {
	tests := []struct {
		name    string
		a, b    any
		want    bool
		wantErr bool
	}{
		{"int greater", 2, 1, false, false},
		{"int equal", 1, 1, false, false},
		{"int less", 1, 2, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LogicLt(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogicLt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LogicLt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogicLte(t *testing.T) {
	tests := []struct {
		name    string
		a, b    any
		want    bool
		wantErr bool
	}{
		{"int greater", 2, 1, false, false},
		{"int equal", 1, 1, true, false},
		{"int less", 1, 2, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LogicLte(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogicLte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LogicLte() = %v, want %v", got, tt.want)
			}
		})
	}
}

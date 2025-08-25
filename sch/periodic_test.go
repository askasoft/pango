package sch

import (
	"testing"
)

func TestPeriodicCron(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want string
	}{
		{"daily", "d 0 10", "0 10 * * *"},
		{"weekly", "w 3 8", "0 8 * * 3"},
		{"monthly", "m 15 22", "0 22 15 * *"},
		{"monthly last day", "m 32 22", "0 22 32 * *"},
		{"invalid", "x 0 0", "x 0 0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := ParsePeriodic(tt.expr)
			if p.Cron() != tt.want {
				t.Fatalf("Cron() = %q, want %q", p.Cron(), tt.want)
			}
			if err == nil {
				if _, err := ParseCron(p.Cron()); err != nil {
					t.Fatalf("ParseCron() = %v, want nil", err)
				}
			}
		})
	}
}

func TestPeriodicString(t *testing.T) {
	p := Periodic{"w 5 14"}
	want := "w 5 14"
	if got := p.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestParsePeriodic(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr bool
	}{
		{
			name: "valid daily",
			expr: "d 0 12,13",
		},
		{
			name: "valid weekly",
			expr: "w 3,4 8-20",
		},
		{
			name: "valid monthly",
			expr: "m 15-20 23,1-3",
		},
		{
			name: "valid monthly last day",
			expr: "m 32 23",
		},
		{
			name:    "invalid field count",
			expr:    "d 12",
			wantErr: true,
		},
		{
			name:    "invalid hour",
			expr:    "d 0 25",
			wantErr: true,
		},
		{
			name:    "invalid weekly day",
			expr:    "w 8 10",
			wantErr: true,
		},
		{
			name:    "invalid monthly day",
			expr:    "m 33 5",
			wantErr: true,
		},
		{
			name:    "invalid unit",
			expr:    "x 1 2",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePeriodic(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParsePeriodic(%q) = %v, wantErr %v", tt.expr, err, tt.wantErr)
			}
			if !tt.wantErr && got.expression != tt.expr {
				t.Errorf("ParsePeriodic() = %q, want %q", got.expression, tt.expr)
			}
		})
	}
}

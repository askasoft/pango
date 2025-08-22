package sch

import (
	"testing"
)

func TestPeriodicCron(t *testing.T) {
	tests := []struct {
		name string
		p    Periodic
		want string
	}{
		{"daily", Periodic{Unit: Daily, Day: 0, Hour: 10}, "0 10 * * *"},
		{"weekly", Periodic{Unit: Weekly, Day: 3, Hour: 8}, "0 8 * * 3"},
		{"monthly", Periodic{Unit: Monthly, Day: 15, Hour: 22}, "0 22 15 * *"},
		{"monthly last day", Periodic{Unit: Monthly, Day: 32, Hour: 22}, "0 22 32 * *"},
		{"invalid", Periodic{Unit: 'x', Day: 0, Hour: 0}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cron := tt.p.Cron()
			if cron != tt.want {
				t.Fatalf("Cron() = %q, want %q", cron, tt.want)
			}
			if cron != "" {
				if _, err := ParseCron(cron); err != nil {
					t.Fatalf("ParseCron() = %v, want nil", err)
				}
			}
		})
	}
}

func TestPeriodicString(t *testing.T) {
	p := Periodic{Unit: Weekly, Day: 5, Hour: 14}
	want := "w 5 14"
	if got := p.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestParsePeriodic(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    Periodic
		wantErr bool
	}{
		{
			name: "valid daily",
			expr: "d 0 12",
			want: Periodic{Unit: Daily, Day: 0, Hour: 12},
		},
		{
			name: "valid weekly",
			expr: "w 3 8",
			want: Periodic{Unit: Weekly, Day: 3, Hour: 8},
		},
		{
			name: "valid monthly",
			expr: "m 15 23",
			want: Periodic{Unit: Monthly, Day: 15, Hour: 23},
		},
		{
			name: "valid monthly last day",
			expr: "m 32 23",
			want: Periodic{Unit: Monthly, Day: 32, Hour: 23},
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
				t.Fatalf("ParsePeriodic() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParsePeriodic() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

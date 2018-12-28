package decconv

import (
	"testing"
)

func TestDecode32(t *testing.T) {
	type args struct {
		precision int
		scale     int
		input     string
	}
	tests := []struct {
		name    string
		args    args
		want    uint32
		conv    string
		wantErr bool
	}{
		{
			name: "only-integral-positive",
			args: args{
				precision: 9,
				scale:     0,
				input:     "12",
			},
			want:    12,
			conv:    "12",
			wantErr: false,
		},
		{
			name: "only-integral-negative",
			args: args{
				precision: 9,
				scale:     0,
				input:     "-123",
			},
			want:    ^uint32(123) + 1,
			conv:    "-123",
			wantErr: false,
		},
		{
			name: "generic-empty-fraction",
			args: args{
				precision: 9,
				scale:     1,
				input:     "12.0",
			},
			want:    120,
			conv:    "12",
			wantErr: false,
		},
		{
			name: "generic-empty-fraction",
			args: args{
				precision: 9,
				scale:     2,
				input:     "12.02",
			},
			want:    1202,
			conv:    "12.02",
			wantErr: false,
		},
		{
			name: "check-input-positive",
			args: args{
				precision: 9,
				scale:     5,
				input:     "3015.07654",
			},
			want:    301507654,
			conv:    "3015.07654",
			wantErr: false,
		},
		{
			name: "check-input-negative",
			args: args{
				precision: 9,
				scale:     5,
				input:     "-3015.07654",
			},
			want:    ^uint32(301507654) + 1,
			conv:    "-3015.07654",
			wantErr: false,
		},
		{
			name: "passing-leading-zeroes",
			args: args{
				precision: 9,
				scale:     2,
				input:     "0000123.25",
			},
			want:    12325,
			conv:    "123.25",
			wantErr: false,
		},
		{
			name: "passing-trailing-zeroes",
			args: args{
				precision: 9,
				scale:     3,
				input:     "123.02500000000000",
			},
			want:    123025,
			conv:    "123.025",
			wantErr: false,
		},
		{
			name: "passing-both-sides-zeroes",
			args: args{
				precision: 9,
				scale:     4,
				input:     "-0000000123.12300000000",
			},
			want:    ^uint32(1231230) + 1,
			conv:    "-123.123",
			wantErr: false,
		},
		{
			name: "error-empty-input",
			args: args{
				precision: 9,
				scale:     4,
				input:     "",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error-invalid-input-integral-part",
			args: args{
				precision: 9,
				scale:     4,
				input:     "1A2.3",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error-invalid-input-fraction-part",
			args: args{
				precision: 9,
				scale:     4,
				input:     "12.b3",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error-overflow-integral",
			args: args{
				precision: 9,
				scale:     6,
				input:     "1234.12",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error-overflow-fraction",
			args: args{
				precision: 9,
				scale:     3,
				input:     "1234.1234",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode32(tt.args.precision, tt.args.scale, []byte(tt.args.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decode32() = %v, want %v", got, tt.want)
			}
			if conv := Encode32(tt.args.scale, got); conv != tt.conv && err == nil {
				t.Errorf("Encode32() = %v, want %v", conv, tt.conv)
			}
		})
	}
}

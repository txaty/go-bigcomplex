package complex

import (
	"math/big"
	"reflect"
	"testing"
)

func TestGaussianInt_Conj(t *testing.T) {
	type fields struct {
		R *big.Int
		I *big.Int
	}
	type args struct {
		origin *GaussianInt
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *GaussianInt
	}{
		{
			name: "test1",
			fields: fields{
				R: big.NewInt(1),
				I: big.NewInt(1),
			},
			args: args{
				origin: NewGaussianInt(big.NewInt(1), big.NewInt(1)),
			},
			want: NewGaussianInt(big.NewInt(1), big.NewInt(-1)),
		},
		{
			name: "test2",
			fields: fields{
				R: big.NewInt(1),
				I: big.NewInt(1),
			},
			args: args{
				origin: NewGaussianInt(big.NewInt(1), big.NewInt(2)),
			},
			want: NewGaussianInt(big.NewInt(1), big.NewInt(-2)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GaussianInt{
				R: tt.fields.R,
				I: tt.fields.I,
			}
			if g.Conj(tt.args.origin); !reflect.DeepEqual(g, tt.want) {
				t.Errorf("Conj() = %v, want %v", g, tt.want)
			}
		})
	}
}

func TestGaussianInt_Div(t *testing.T) {
	type args struct {
		a *GaussianInt
		b *GaussianInt
	}
	tests := []struct {
		name         string
		args         args
		wantReminder *GaussianInt
		wantQuotient *GaussianInt
	}{
		{
			name: "test_(1,1)_(1,1)",
			args: args{
				a: NewGaussianInt(big.NewInt(1), big.NewInt(1)),
				b: NewGaussianInt(big.NewInt(1), big.NewInt(1)),
			},
			wantReminder: NewGaussianInt(big.NewInt(0), big.NewInt(0)),
			wantQuotient: NewGaussianInt(big.NewInt(1), big.NewInt(0)),
		},
		{
			name: "test_(1,1)_(2,2)",
			args: args{
				a: NewGaussianInt(big.NewInt(1), big.NewInt(1)),
				b: NewGaussianInt(big.NewInt(2), big.NewInt(2)),
			},
			wantReminder: NewGaussianInt(big.NewInt(1), big.NewInt(1)),
			wantQuotient: NewGaussianInt(big.NewInt(0), big.NewInt(0)),
		},
		{
			name: "test_(7,3)_(2,-1)",
			args: args{
				a: NewGaussianInt(big.NewInt(7), big.NewInt(3)),
				b: NewGaussianInt(big.NewInt(2), big.NewInt(-1)),
			},
			wantReminder: NewGaussianInt(big.NewInt(0), big.NewInt(-1)),
			wantQuotient: NewGaussianInt(big.NewInt(2), big.NewInt(3)),
		},
		{
			name: "test_(2,1)_(1,0)",
			args: args{
				a: NewGaussianInt(big.NewInt(2), big.NewInt(1)),
				b: NewGaussianInt(big.NewInt(1), big.NewInt(0)),
			},
			wantReminder: NewGaussianInt(big.NewInt(0), big.NewInt(0)),
			wantQuotient: NewGaussianInt(big.NewInt(2), big.NewInt(1)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GaussianInt{}
			quotient := g.Div(tt.args.a, tt.args.b)
			if g.R.Cmp(tt.wantReminder.R) != 0 || g.I.Cmp(tt.wantReminder.I) != 0 {
				t.Errorf("g = %v, want reminder %v", g, tt.wantReminder)
			}
			if quotient.R.Cmp(tt.wantQuotient.R) != 0 || quotient.I.Cmp(tt.wantQuotient.I) != 0 {
				t.Errorf("Div() = %v, want quotient %v", quotient, tt.wantQuotient)
			}
		})
	}
}

func TestGaussianInt_String(t *testing.T) {
	type fields struct {
		R *big.Int
		I *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test_1+i",
			fields: fields{
				R: big.NewInt(1),
				I: big.NewInt(1),
			},
			want: "1+i",
		},
		{
			name: "test_1-i",
			fields: fields{
				R: big.NewInt(1),
				I: big.NewInt(-1),
			},
			want: "1-i",
		},
		{
			name: "test_-1+i",
			fields: fields{
				R: big.NewInt(-1),
				I: big.NewInt(1),
			},
			want: "-1+i",
		},
		{
			name: "test_-1-i",
			fields: fields{
				R: big.NewInt(-1),
				I: big.NewInt(-1),
			},
			want: "-1-i",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GaussianInt{
				R: tt.fields.R,
				I: tt.fields.I,
			}
			if got := g.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

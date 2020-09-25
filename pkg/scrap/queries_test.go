package scrap

import (
	"strings"
	"testing"
)

func Test_buildQuery(t *testing.T) {
	type args struct {
		provider string
		id       string
		region   string
	}
	tests := []struct {
		name    string
		args    args
		wantURL string
		wantErr bool
	}{
		{
			"fails because provided not supported",
			args{"unsupported", "B00JKEJ4TA", "de"},
			"",
			true,
		},
		{
			"succeeds since all is good",
			args{AmazonPrime, "B00JKEJ4TA", "de"},
			"https://www.amazon.de/gp/product/B00JKEJ4TA",
			false,
		},
		{
			"succeeds since all is good",
			args{AmazonPrime, "B00JKEJ4TA", "com"},
			"https://www.amazon.com/gp/product/B00JKEJ4TA",
			false,
		},
		{
			"fails because the id is too short",
			args{AmazonPrime, "SHORT", "de"},
			"",
			true,
		},
		{
			"fails because the id is too long",
			args{AmazonPrime, "WAYTOOLONGFORANID", "de"},
			"",
			true,
		},
		{
			"fails because the id has wrong case",
			args{AmazonPrime, "wrongcase", "de"},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQ, err := buildQuery(tt.args.provider, tt.args.id, tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !strings.Contains(gotQ, tt.wantURL) {
				t.Errorf("buildQuery()  query = %v, want %v in query", gotQ, tt.wantURL)
			}
		})
	}
}

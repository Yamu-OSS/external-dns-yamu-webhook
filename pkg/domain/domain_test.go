package domain

import "testing"

func TestHasSuffix(t *testing.T) {
	type args struct {
		s      string
		suffix string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				s:      "abc.COM",
				suffix: "com",
			},
			want: true,
		},
		{
			name: "test2",
			args: args{
				s:      "abc.CoM.",
				suffix: "com.",
			},
			want: true,
		},
		{
			name: "test3",
			args: args{
				s:      "abc.CON",
				suffix: "com",
			},
			want: false,
		},
		{
			name: "test4",
			args: args{
				s:      "abc.com.",
				suffix: ".",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasSuffix(tt.args.s, tt.args.suffix); got != tt.want {
				t.Errorf("HasSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrimSuffix(t *testing.T) {
	type args struct {
		s      string
		suffix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				s:      "abc.COM",
				suffix: "com",
			},
			want: "abc.",
		},
		{
			name: "test2",
			args: args{
				s:      "abc.CoM.",
				suffix: "com.",
			},
			want: "abc.",
		},
		{
			name: "test3",
			args: args{
				s:      "abc.CON",
				suffix: "com",
			},
			want: "abc.CON",
		},
		{
			name: "test4",
			args: args{
				s:      "abc.com.",
				suffix: ".",
			},
			want: "abc.com",
		},
		{
			name: "test5",
			args: args{
				s:      "ABC.com.",
				suffix: "COM.",
			},
			want: "ABC.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimSuffix(tt.args.s, tt.args.suffix); got != tt.want {
				t.Errorf("TrimSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitSuffixToDomain(t *testing.T) {
	type args struct {
		s       string
		suffixs []string
	}
	tests := []struct {
		name     string
		args     args
		wantPre  string
		wantSuff string
	}{
		{
			name: "test1",
			args: args{
				s:       "abc.COM",
				suffixs: []string{"com"},
			},
			wantPre:  "abc",
			wantSuff: "com",
		},
		{
			name: "test2",
			args: args{
				s:       "abc.COM",
				suffixs: []string{"com."},
			},
			wantPre:  "abc",
			wantSuff: "com",
		},
		{
			name: "test3",
			args: args{
				s:       "abc.CoM.",
				suffixs: []string{"com", "c.COM"},
			},
			wantPre:  "abc",
			wantSuff: "com",
		},
		{
			name: "test4",
			args: args{
				s:       "abc.CoM.",
				suffixs: []string{"com", "abc.COM"},
			},
			wantPre:  "",
			wantSuff: "abc.COM",
		},
		{
			name: "test5",
			args: args{
				s:       "abc.CON.",
				suffixs: []string{"com"},
			},
			wantPre:  "abc.CON",
			wantSuff: "", // 未匹配到
		},
		{
			name: "test5-1",
			args: args{
				s:       "abc.COM.",
				suffixs: []string{"e.abc.com"},
			},
			wantPre:  "abc.COM",
			wantSuff: "", // 未匹配到
		},
		{
			name: "test6",
			args: args{
				s:       "abc.com.",
				suffixs: []string{"."},
			},
			wantPre:  "abc.com",
			wantSuff: ".",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre, suff := SplitSuffixToDomain(tt.args.s, tt.args.suffixs)
			if pre != tt.wantPre || suff != tt.wantSuff {
				t.Errorf("SplitSuffixToDomain() = %v; %v, want %v; %v", pre, suff, tt.wantPre, tt.wantSuff)
			}
		})
	}
}

func TestHostAddDomain(t *testing.T) {
	type args struct {
		host   string
		domain string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				host:   "www",
				domain: "test.com",
			},
			want: "www.test.com",
		},
		{
			name: "test2",
			args: args{
				host:   "www",
				domain: "test.com.",
			},
			want: "www.test.com",
		},
		{
			name: "test3",
			args: args{
				host:   "abc",
				domain: ".",
			},
			want: "abc",
		},
		{
			name: "test4",
			args: args{
				host:   "",
				domain: ".",
			},
			want: "",
		},
		{
			name: "test5",
			args: args{
				host:   "",
				domain: "abc.",
			},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HostAddDomain(tt.args.host, tt.args.domain); got != tt.want {
				t.Errorf("HostAddDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

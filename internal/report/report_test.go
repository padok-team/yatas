package report

import (
	"reflect"
	"testing"

	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func Test_countResultOkOverall(t *testing.T) {
	type args struct {
		results []results.Result
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
	}{
		{
			name: "count result ok overall",
			args: args{
				results: []results.Result{
					{
						Status: "OK",
					},
					{
						Status: "OK",
					},
				},
			},
			want:  2,
			want1: 2,
		},
		{
			name: "count result ok overall",
			args: args{
				results: []results.Result{
					{
						Status: "FAIL",
					},
					{
						Status: "OK",
					},
				},
			},
			want:  1,
			want1: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := countResultOkOverall(tt.args.results)
			if got != tt.want {
				t.Errorf("countResultOkOverall() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("countResultOkOverall() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestIsIgnored(t *testing.T) {
	type args struct {
		c     *yatas.Config
		r     results.Result
		check results.Check
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "is ignored",
			args: args{
				c: &yatas.Config{
					Ignore: []yatas.Ignore{
						{
							ID:    "test",
							Regex: true,
							Values: []string{
								"test",
							},
						},
					},
				},
				r: results.Result{
					Message:    "test",
					Status:     "OK",
					ResourceID: "test",
				},
				check: results.Check{
					Id: "test",
				},
			},
			want: true,
		},
		{
			name: "is ignored",
			args: args{
				c: &yatas.Config{
					Ignore: []yatas.Ignore{
						{
							ID:    "test",
							Regex: false,
							Values: []string{
								"test",
							},
						},
					},
				},
				r: results.Result{
					Message:    "test",
					Status:     "OK",
					ResourceID: "test",
				},
				check: results.Check{
					Id: "test",
				},
			},
			want: true,
		},
		{
			name: "is ignored",
			args: args{
				c: &yatas.Config{
					Ignore: []yatas.Ignore{
						{
							ID:    "test",
							Regex: false,
							Values: []string{
								"test",
							},
						},
					},
				},
				r: results.Result{
					Message:    "test",
					Status:     "OK",
					ResourceID: "test",
				},
				check: results.Check{
					Id: "toto",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIgnored(tt.args.c, tt.args.r, tt.args.check); got != tt.want {
				t.Errorf("IsIgnored() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExitCode(t *testing.T) {
	type args struct {
		checks []results.Tests
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "exit code",
			args: args{
				checks: []results.Tests{
					{
						Account: "test",
						Checks: []results.Check{
							{
								Id: "test",
								Results: []results.Result{
									{
										Status: "OK",
									},
								},
								Status: "OK",
							},
						},
					},
				},
			},
			want: 0,
		},
		{
			name: "exit code",
			args: args{
				checks: []results.Tests{
					{
						Account: "test",
						Checks: []results.Check{
							{
								Id: "test",
								Results: []results.Result{
									{
										Status: "FAIL",
									},
								},
								Status: "FAIL",
							},
						},
					},
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExitCode(tt.args.checks); got != tt.want {
				t.Errorf("ExitCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountChecksPassedOverall(t *testing.T) {
	type args struct {
		checks []results.Check
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
	}{
		{
			name: "count checks passed overall",
			args: args{
				checks: []results.Check{
					{
						Id:     "test",
						Status: "OK",
					},
				},
			},
			want:  1,
			want1: 1,
		},
		{
			name: "count checks passed overall",
			args: args{
				checks: []results.Check{
					{
						Id:     "test",
						Status: "OK",
					},
					{
						Id:     "test",
						Status: "FAIL",
					},
				},
			},
			want:  1,
			want1: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CountChecksPassedOverall(tt.args.checks)
			if got != tt.want {
				t.Errorf("CountChecksPassedOverall() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CountChecksPassedOverall() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRemoveIgnored(t *testing.T) {
	type args struct {
		c     *yatas.Config
		tests []results.Tests
	}
	tests := []struct {
		name string
		args args
		want []results.Tests
	}{
		{
			name: "remove ignored",
			args: args{
				c: &yatas.Config{
					Ignore: []yatas.Ignore{
						{
							ID:    "test",
							Regex: true,
							Values: []string{
								"test",
							},
						},
					},
				},
				tests: []results.Tests{
					{
						Account: "test",
						Checks: []results.Check{
							{
								Id: "test",

								Results: []results.Result{
									{
										Status:     "FAIL",
										Message:    "test",
										ResourceID: "test",
									},
									{
										Status:     "OK",
										Message:    "toto",
										ResourceID: "toto",
									},
								},
							},
						},
					},
				},
			},
			want: []results.Tests{
				{
					Account: "test",
					Checks: []results.Check{
						{
							Status: "OK",
							Id:     "test",
							Results: []results.Result{
								{
									Status:     "OK",
									Message:    "toto",
									ResourceID: "toto",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveIgnored(tt.args.c, tt.args.tests); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveIgnored() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

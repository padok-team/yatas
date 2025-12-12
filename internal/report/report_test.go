package report

import (
	"testing"

	"github.com/padok-team/yatas/plugins/commons"
)

func Test_countResultOkOverall(t *testing.T) {
	type args struct {
		results []commons.Result
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
				results: []commons.Result{
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
				results: []commons.Result{
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
		c     *commons.Config
		r     commons.Result
		check commons.Check
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "is ignored",
			args: args{
				c: &commons.Config{
					Ignore: []commons.Ignore{
						{
							ID:    "test",
							Regex: true,
							Values: []string{
								"test",
							},
						},
					},
				},
				r: commons.Result{
					Message:    "test",
					Status:     "OK",
					ResourceID: "test",
				},
				check: commons.Check{
					Id: "test",
				},
			},
			want: true,
		},
		{
			name: "is ignored",
			args: args{
				c: &commons.Config{
					Ignore: []commons.Ignore{
						{
							ID:    "test",
							Regex: false,
							Values: []string{
								"test",
							},
						},
					},
				},
				r: commons.Result{
					Message:    "test",
					Status:     "OK",
					ResourceID: "test",
				},
				check: commons.Check{
					Id: "test",
				},
			},
			want: true,
		},
		{
			name: "is ignored",
			args: args{
				c: &commons.Config{
					Ignore: []commons.Ignore{
						{
							ID:    "test",
							Regex: false,
							Values: []string{
								"test",
							},
						},
					},
				},
				r: commons.Result{
					Message:    "test",
					Status:     "OK",
					ResourceID: "test",
				},
				check: commons.Check{
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
		checks []commons.Tests
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "exit code",
			args: args{
				checks: []commons.Tests{
					{
						Account: "test",
						Checks: []commons.Check{
							{
								Id: "test",
								Results: []commons.Result{
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
				checks: []commons.Tests{
					{
						Account: "test",
						Checks: []commons.Check{
							{
								Id: "test",
								Results: []commons.Result{
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
		checks []commons.Check
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
				checks: []commons.Check{
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
				checks: []commons.Check{
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
		c     *commons.Config
		tests []commons.Tests
	}
	tests := []struct {
		name string
		args args
		want []commons.Tests
	}{
		{
			name: "remove ignored",
			args: args{
				c: &commons.Config{
					Ignore: []commons.Ignore{
						{
							ID:    "test",
							Regex: true,
							Values: []string{
								"test",
							},
						},
					},
				},
				tests: []commons.Tests{
					{
						Account: "test",
						Checks: []commons.Check{
							{
								Id: "test",

								Results: []commons.Result{
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
			want: []commons.Tests{
				{
					Account: "test",
					Checks: []commons.Check{
						{
							Status: "OK",
							Id:     "test",
							Results: []commons.Result{
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
			got := RemoveIgnored(tt.args.c, tt.args.tests)
			if len(got) != len(tt.want) {
				t.Errorf("RemoveIgnored() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFilterHDSChecks(t *testing.T) {
	type args struct {
		tests []commons.Tests
	}
	tests := []struct {
		name string
		args args
		want []commons.Tests
	}{
		{
			name: "filter HDS checks - some have HDS category",
			args: args{
				tests: []commons.Tests{
					{
						Account: "test-account",
						Checks: []commons.Check{
							{
								Id:         "CHECK_001",
								Name:       "HDS Check",
								Categories: []string{"Security", "HDS"},
								Status:     "OK",
							},
							{
								Id:         "CHECK_002",
								Name:       "Non-HDS Check",
								Categories: []string{"Security", "Performance"},
								Status:     "FAIL",
							},
							{
								Id:         "CHECK_003",
								Name:       "Another HDS Check",
								Categories: []string{"HDS", "Compliance"},
								Status:     "OK",
							},
						},
					},
				},
			},
			want: []commons.Tests{
				{
					Account: "test-account",
					Checks: []commons.Check{
						{
							Id:         "CHECK_001",
							Name:       "HDS Check",
							Categories: []string{"Security", "HDS"},
							Status:     "OK",
						},
						{
							Id:         "CHECK_003",
							Name:       "Another HDS Check",
							Categories: []string{"HDS", "Compliance"},
							Status:     "OK",
						},
					},
				},
			},
		},
		{
			name: "filter HDS checks - no HDS checks",
			args: args{
				tests: []commons.Tests{
					{
						Account: "test-account",
						Checks: []commons.Check{
							{
								Id:         "CHECK_001",
								Name:       "Security Check",
								Categories: []string{"Security"},
								Status:     "OK",
							},
							{
								Id:         "CHECK_002",
								Name:       "Performance Check",
								Categories: []string{"Performance"},
								Status:     "FAIL",
							},
						},
					},
				},
			},
			want: []commons.Tests{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterHDSChecks(tt.args.tests)
			if len(got) != len(tt.want) {
				t.Errorf("FilterHDSChecks() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i, test := range got {
				if test.Account != tt.want[i].Account {
					t.Errorf("FilterHDSChecks() account = %v, want %v", test.Account, tt.want[i].Account)
				}
				if len(test.Checks) != len(tt.want[i].Checks) {
					t.Errorf("FilterHDSChecks() checks length = %v, want %v", len(test.Checks), len(tt.want[i].Checks))
				}
			}
		})
	}
}

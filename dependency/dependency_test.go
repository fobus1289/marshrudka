package dependency

import (
	"testing"
)

func Test_dependency_Set_Get(t *testing.T) {

	type args struct {
		s int
	}

	tests := []struct {
		name string
		d    IDependency
		args *args
		want IDependency
	}{
		{
			name: "set",
			d:    NewDependency().Set(&args{s: 1}),
			args: &args{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.d.Get(&tt.args)

			if tt.args.s != 1 {
				t.Errorf("dependency.Set() = %v", tt.args)
			}

		})
	}

}

func Test_dependency_Fill_SetAll(t *testing.T) {

	type T struct {
		name string
	}

	type T1 struct {
		name string
	}

	type T2 struct {
		name string
	}

	type args struct {
		*T
		*T1
		*T2
	}

	tests := []struct {
		name string
		d    IDependency
		args *args
	}{
		{
			name: "fill",
			d:    NewDependency().SetAll(&T{name: "T"}, &T1{name: "T1"}, &T2{name: "T2"}),
			args: &args{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.Fill(tt.args)
			args := tt.args
			if (args.T == nil || args.T1 == nil || args.T2 == nil) ||
				(args.T.name != "T" || args.T1.name != "T1" || args.T2.name != "T2") {
				t.Errorf("dependency.Fill() = %v", args)
			}
		})
	}
}

// 9.720 ns/op
func BenchmarkCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Count(i)
	}
}

// 146.4 ns/op
func BenchmarkSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Split(i)
	}
}

func BenchmarkRegexp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		args := Regexp("api/users/1/user1/edit")
		_ = args
	}
}

func BenchmarkRegexpFindAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		arrs := RegexpFindAll("api/users/1/user1/edit")
		_ = arrs
	}
}

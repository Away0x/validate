package validate_test

import (
	"testing"

	"github.com/Away0x/validate"
)

type User struct {
	validate.BaseValidate
	Name string
	Age  int
}

type user2 User

func (u *user2) Validators() validate.Validators {
	return validate.Validators{
		"name": {
			validate.RequiredValidator(u.Name),
			validate.MinLengthValidator(u.Name, 3),
		},
		"age": {
			func() string {
				if u.Age < 10 {
					return "age 不能小于 10"
				}
				return ""
			},
		},
	}
}

type user3 User

func (*user3) IsStrict() bool {
	return false
}
func (u *user3) Validators() validate.Validators {
	return validate.Validators{
		"name": {
			validate.RequiredValidator(u.Name),
			validate.MinLengthValidator(u.Name, 3),
		},
		"age": {
			func() string {
				if u.Age < 10 {
					return "age 不能小于 10"
				}
				return ""
			},
		},
	}
}

type user4 User

func (u *user4) Validators() validate.Validators {
	return validate.Validators{
		"name": {
			validate.RequiredValidator(u.Name),
		},
	}
}
func (u *user4) Messages() validate.Messages {
	return validate.Messages{
		"name": {
			"用户名必须存在",
		},
	}
}

type user5 User

func (u *user5) Plugins() validate.Plugins {
	return validate.Plugins{
		func() (string, []validate.ValidatorFunc, []string) {
			return "name", []validate.ValidatorFunc{validate.RequiredValidator(u.Name)}, []string{"用户名必须存在"}
		},
	}
}

func TestValidate(t *testing.T) {
	// --------------
	u := &User{
		Name: "xiaoming",
		Age:  18,
	}

	if _, ok := validate.Run(u); !ok {
		t.Error("u validate error")
	}

	// --------------
	u2 := &user2{
		Name: "",
		Age:  7,
	}

	if msg, ok := validate.Run(u2); !ok {
		if m, ok2 := msg["age"]; ok2 {
			t.Error("u2 age msg error " + m[0])
		}
	} else {
		t.Error("u2 validate error")
	}

	// --------------
	u3 := &user3{
		Name: "",
		Age:  7,
	}

	if msg, ok := validate.Run(u3); !ok {
		if m, ok2 := msg["name"]; !ok2 {
			t.Error("u3 name msg error")
		} else {
			if len(m) != 2 || m[0] != "name 必须存在" || m[1] != "name 必须大于 3 个字符" {
				t.Error("u3 name msg error")
			}
		}

		if m, ok2 := msg["age"]; !ok2 {
			t.Error("u3 age msg error")
		} else {
			if len(m) != 1 || m[0] != "age 不能小于 10" {
				t.Error("u3 age msg error")
			}
		}
	} else {
		t.Error("u3 validate error")
	}

	// --------------
	u31 := &user3{
		Name: "abcd",
		Age:  11,
	}
	if _, ok := validate.Run(u31); !ok {
		t.Error("u31 validate error")
	}

	// --------------
	u4 := &user4{
		Name: "",
	}
	if msg, ok := validate.Run(u4); !ok {
		if m, ok2 := msg["name"]; !ok2 {
			t.Error("u4 name msg error")
		} else {
			if m[0] != "用户名必须存在" {
				t.Error("u4 name msg error " + m[0])
			}
		}
	} else {
		t.Error("u4 validate error")
	}

	// --------------
	u5 := &user5{
		Name: "",
	}
	if msg, ok := validate.Run(u5); !ok {
		if m, ok2 := msg["name"]; !ok2 {
			t.Error("u5 name msg error")
		} else {
			if m[0] != "用户名必须存在" {
				t.Error("u5 name msg error " + m[0])
			}
		}
	} else {
		t.Error("u5 validate error")
	}

	// --------------
	u6 := &User{Name: ""}
	if msg, ok := validate.RunWithConfig(u6, validate.Config{
		Plugins: validate.Plugins{
			func() (string, []validate.ValidatorFunc, []string) {
				return "name", []validate.ValidatorFunc{validate.RequiredValidator(u6.Name)}, []string{"用户名必须存在"}
			},
		},
	}); !ok {
		if m, ok2 := msg["name"]; !ok2 {
			t.Error("u6 name msg error")
		} else {
			if m[0] != "用户名必须存在" {
				t.Error("u6 name msg error " + m[0])
			}
		}
	} else {
		t.Error("u6 validate error")
	}
}

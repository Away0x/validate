# Validate
简单易用的 go struct 字段验证插件

> go get -u github.com/Away0x/validate

## example

```go
type loginForm struct {
	validate.BaseValidate
	Email    string
	Password string
}

func (f *loginForm) Validators() validate.Validators {
	return validate.Validators{
		"email": {
			validate.RequiredValidator(f.Email),
			validate.EmailValidator(f.Email),
			validate.MaxLengthValidator(f.Email, 255),
		},
		"password": {
			validate.RequiredValidator(f.Password),
		},
	}
}

func (f *loginForm) Messages() validate.Messages {
	return validate.Messages{
		"email": {
			"邮箱不能为空",
			"邮箱格式错误",
			"邮箱长度不能大于 255 个字符",
		},
		"password": {
			"密码不能为空",
		},
	}
}

func (*loginForm) IsStrict() bool {
	return false
}

form := &loginForm{
  Email: "",
  Password: "",
}

if errMap, ok := validate.Run(form); !ok {
  b, _ := json.MarshalIndent(errMap, "", "\t")
  fmt.Println(string(b))
}

/*
{
  "email": [
    "邮箱不能为空"
  ],
  "password": [
    "密码不能为空"
  ]
}
*/
```

## plugin example
```go
// EmailPlugin 项目中 email 字段请求参数的验证
func EmailPlugin(email string) validate.PluginFunc {
	return func() (string, []validate.ValidatorFunc, []string) {
		return "email", []validate.ValidatorFunc{
				validate.RequiredValidator(email),
				validate.EmailValidator(email),
				validate.MaxLengthValidator(email, 255),
			}, []string{
				"邮箱不能为空",
				"邮箱格式错误",
				"邮箱长度不能大于 255 个字符",
			}
	}
}

// PasswordPlugin 项目中 password 字段请求参数的验证
func PasswordPlugin(password string) validate.PluginFunc {
	return func() (string, []validate.ValidatorFunc, []string) {
		return "password", []validate.ValidatorFunc{
				validate.RequiredValidator(password),
			}, []string{
				"密码不能为空",
			}
	}
}
```
```go
type loginForm struct {
	validate.BaseValidate
	Email    string
	Password string
}

func (*loginForm) IsStrict() bool {
	return false
}

func (f *loginForm) Plugins() []validate.PluginFunc {
	return []validate.PluginFunc{
		EmailPlugin(f.Email),
		PasswordPlugin(f.Password),
	}
}

form := &loginForm{
  Email: "",
  Password: "",
}

if errMap, ok := validate.Run(form); !ok {
  b, _ := json.MarshalIndent(errMap, "", "\t")
  fmt.Println(string(b))
}

/*
{
  "email": [
    "邮箱不能为空"
  ],
  "password": [
    "密码不能为空"
  ]
}
*/
```

## withConfig example
```go
type loginForm struct {
	validate.BaseValidate
	Email    string
	Password string
}

form := &loginForm{
  Email: "",
  Password: "",
}

errMap, ok := validate.RunWithConfig(form, validate.Config{
  Plugins: validate.Plugins{
    EmailPlugin(req.Email),
    PasswordPlugin(req.Password),
  },
})
if !ok {
  b, _ := json.MarshalIndent(errMap, "", "\t")
  fmt.Println(string(b))
}

/*
{
  "email": [
    "邮箱不能为空"
  ],
  "password": [
    "密码不能为空"
  ]
}
*/
```

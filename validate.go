package validate

import (
	"strings"
)

type (
	// ValidatorFunc 验证器函数 (返回的 string 不为空即表示验证失败)
	ValidatorFunc func() (msg string)
	// Validators 验证器数组 map
	Validators map[string][]ValidatorFunc
	// Messages 错误信息
	Messages map[string][]string
	// PluginFunc 聚合 Validators 和 Messages 的方法
	PluginFunc func() (string, []ValidatorFunc, []string)
	// Plugins -
	Plugins []PluginFunc

	// Validater -
	Validater interface {
		// IsStrict : 严格模式时，第一个验证出错时，即会停止其他验证
		IsStrict() bool
		// Validators : 注册验证器 map
		Validators() Validators
		// Messages : 注册错误信息 map
		Messages() Messages
		// Plugins : 提供一个聚合验证器和错误信息的方法
		// (Plugins 方法中注册的数据会被 Validators/Messages 方法注册的数据覆盖)
		Plugins() Plugins
	}

	// Config : RunWithConfig 时用到
	Config struct {
		Strict     bool
		Validators Validators
		Messages   Messages
		Plugins    Plugins
	}
)

// Error error interface
func (msg Messages) Error() string {
	var val strings.Builder
	for k, v := range msg {
		val.WriteString(k + ": " + strings.Join(v, ",") + "\n")
	}

	return val.String()
}

// BaseValidate -
type BaseValidate struct{}

// IsStrict 是否为严格模式
func (*BaseValidate) IsStrict() bool {
	return true
}

// Validators : 注册验证器
// 验证器数组按顺序验证，一旦验证没通过，即结束该字段的验证
func (*BaseValidate) Validators() Validators {
	return Validators{}
}

// Messages 注册错误信息
func (*BaseValidate) Messages() Messages {
	return Messages{}
}

// Plugins 注册 plugin
func (*BaseValidate) Plugins() Plugins {
	return nil
}

// Run 执行验证
func Run(v Validater) (errMap Messages, ok bool) {
	var (
		strict       = v.IsStrict()
		plugins      = v.Plugins()
		validatorMap = Validators{}
		messageMap   = Messages{}
	)

	// 1. ------------------------ 收集验证器和错误信息 ------------------------
	// 获取到 #Plugins 里面的 validators 和 messages
	if plugins != nil || len(plugins) != 0 {
		for _, plugin := range plugins {
			key, validators, messages := plugin()
			validatorMap[key] = validators
			messageMap[key] = messages
		}
	}

	// 获取到  #Validators 和 #Messages 方法的数据
	validatorsFuncResult := v.Validators()
	messagesFuncResult := v.Messages()
	for key, validators := range validatorsFuncResult {
		validatorMap[key] = validators
		messageMap[key] = messagesFuncResult[key]
	}

	// 2. ------------------------ 执行验证 ------------------------
	return runValidate(strict, validatorMap, messageMap)
}

// RunWithConfig -
func RunWithConfig(v Validater, config Config) (errMap Messages, ok bool) {
	var (
		strict       = config.Strict
		plugins      = config.Plugins
		validatorMap = Validators{}
		messageMap   = Messages{}
	)

	// 1. ------------------------ 收集验证器和错误信息 ------------------------
	if plugins != nil || len(plugins) != 0 {
		for _, plugin := range plugins {
			key, validators, messages := plugin()
			validatorMap[key] = validators
			messageMap[key] = messages
		}
	}

	validatorsFuncResult := config.Validators
	messagesFuncResult := config.Messages
	if validatorsFuncResult != nil {
		for key, validators := range validatorsFuncResult {
			validatorMap[key] = validators
			if messagesFuncResult != nil {
				messageMap[key] = messagesFuncResult[key]
			}
		}
	}

	// 2. ------------------------ 执行验证 ------------------------
	return runValidate(strict, validatorMap, messageMap)
}

// 执行验证
func runValidate(strict bool, validatorMap Validators, messageMap Messages) (errMap Messages, ok bool) {
	errMap = make(Messages)
	ok = true

	for key, validators := range validatorMap {
		customMsgArr := messageMap[key] // 自定义的错误信息
		customMsgArrLen := len(customMsgArr)

		for i, fn := range validators {
			errMsg := fn() // 执行验证函数
			if errMsg != "" {
				ok = false

				if i < customMsgArrLen && customMsgArr[i] != "" {
					// 采用自定义的错误信息输出
					errMsg = customMsgArr[i]
				} else {
					// 采用默认的错误信息输出，错误信息中可使用 $name 替换字段名
					errMsg = parseEasyTemplate(errMsg, map[string]string{
						"$name": key,
					})
				}

				if errMap[key] == nil {
					errMap[key] = make([]string, 0)
				}
				errMap[key] = append(errMap[key], errMsg)

				if strict {
					return // 严格模式: 结束所有验证
				}
			}
		}
	}

	return
}

// 解析 string 模板
func parseEasyTemplate(tplString string, data map[string]string) string {
	replaceArr := []string{}
	for k, v := range data {
		replaceArr = append(replaceArr, k, v)
	}

	r := strings.NewReplacer(replaceArr...)

	return r.Replace(tplString)
}

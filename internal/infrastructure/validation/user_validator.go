package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/taskflow/internal/domain/service"
)

// UserValidator 用户验证器实现
type UserValidator struct {
	emailRegex    *regexp.Regexp
	weakPasswords map[string]bool
}

// NewUserValidator 创建用户验证器
func NewUserValidator() service.UserValidator {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// 常见弱密码列表
	weakPasswords := map[string]bool{
		"password":    true,
		"123456":      true,
		"12345678":    true,
		"qwerty":      true,
		"abc123":      true,
		"password123": true,
		"admin":       true,
		"root":        true,
		"guest":       true,
		"test":        true,
	}

	return &UserValidator{
		emailRegex:    emailRegex,
		weakPasswords: weakPasswords,
	}
}

// ValidateEmail 验证邮箱格式
func (v *UserValidator) ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("邮箱不能为空")
	}

	if len(email) > 254 {
		return fmt.Errorf("邮箱长度不能超过254个字符")
	}

	if !v.emailRegex.MatchString(email) {
		return fmt.Errorf("邮箱格式不正确")
	}

	// 检查邮箱域名部分
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("邮箱格式不正确")
	}

	domain := parts[1]
	if len(domain) > 253 {
		return fmt.Errorf("邮箱域名过长")
	}

	return nil
}

// ValidatePassword 验证密码强度
func (v *UserValidator) ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("密码不能为空")
	}

	if len(password) < 8 {
		return fmt.Errorf("密码长度至少8个字符")
	}

	if len(password) > 128 {
		return fmt.Errorf("密码长度不能超过128个字符")
	}

	// 检查是否为弱密码
	if v.weakPasswords[strings.ToLower(password)] {
		return fmt.Errorf("密码过于简单，请使用更复杂的密码")
	}

	// 检查密码复杂度
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	var missing []string
	if !hasUpper {
		missing = append(missing, "大写字母")
	}
	if !hasLower {
		missing = append(missing, "小写字母")
	}
	if !hasNumber {
		missing = append(missing, "数字")
	}
	if !hasSpecial {
		missing = append(missing, "特殊字符")
	}

	// 至少需要包含3种类型的字符
	if len(missing) > 1 {
		return fmt.Errorf("密码必须包含以下字符类型中的至少3种：大写字母、小写字母、数字、特殊字符。当前缺少：%s", strings.Join(missing, "、"))
	}

	return nil
}

// ValidateUsername 验证用户名
func (v *UserValidator) ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}

	if len(username) < 3 {
		return fmt.Errorf("用户名长度至少3个字符")
	}

	if len(username) > 50 {
		return fmt.Errorf("用户名长度不能超过50个字符")
	}

	// 用户名只能包含字母、数字、下划线和连字符
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("用户名只能包含字母、数字、下划线和连字符")
	}

	// 不能以数字开头
	if unicode.IsNumber(rune(username[0])) {
		return fmt.Errorf("用户名不能以数字开头")
	}

	return nil
}

// ValidateName 验证姓名
func (v *UserValidator) ValidateName(name string) error {
	return v.ValidateFullName(name)
}

// ValidateFullName 验证全名
func (v *UserValidator) ValidateFullName(fullName string) error {
	if fullName == "" {
		return fmt.Errorf("姓名不能为空")
	}

	if len(fullName) < 2 {
		return fmt.Errorf("姓名长度至少2个字符")
	}

	if len(fullName) > 100 {
		return fmt.Errorf("姓名长度不能超过100个字符")
	}

	// 检查是否包含非法字符
	for _, char := range fullName {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) && char != '·' && char != '-' {
			return fmt.Errorf("姓名包含非法字符")
		}
	}

	return nil
}

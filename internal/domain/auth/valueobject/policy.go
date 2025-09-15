package valueobject

import "encoding/json"

// PolicyRule ABAC策略规则
type PolicyRule struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Resource    string                 `json:"resource"`
	Action      string                 `json:"action"`
	Effect      PolicyEffect           `json:"effect"`     // allow, deny
	Conditions  map[string]interface{} `json:"conditions"` // 条件表达式
	Priority    int                    `json:"priority"`   // 优先级，数值越大优先级越高
	IsActive    bool                   `json:"is_active"`
}

// PolicyID 策略ID值对象
type PolicyID string

func (id PolicyID) String() string {
	return string(id)
}

// PolicyEffect 策略效果值对象
type PolicyEffect string

const (
	PolicyEffectAllow PolicyEffect = "allow"
	PolicyEffectDeny  PolicyEffect = "deny"
)

// PolicyConditions 策略条件值对象
type PolicyConditions map[string]interface{}

// ToJSON 将条件转换为JSON字符串
func (pc PolicyConditions) ToJSON() (string, error) {
	data, err := json.Marshal(pc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON 从JSON字符串解析条件
func (pc *PolicyConditions) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), pc)
}

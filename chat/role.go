package chat

type Role string

const (
	RoleAssistant Role = "assistant"
	RoleDeveloper Role = "developer"
	RoleSystem    Role = "system"
	RoleTool      Role = "tool"
	RoleUser      Role = "user"
	RoleModel     Role = "model"
)

func (r Role) String() string {
	return string(r)
}

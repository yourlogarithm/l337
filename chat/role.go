package chat

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
	RoleDeveloper Role = "developer"
	RoleFunction  Role = "function"
)

func (r Role) String() string {
	return string(r)
}

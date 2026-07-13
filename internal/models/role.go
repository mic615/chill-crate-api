package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Role int

const (
	RoleViewer Role = iota // 0 (least privilege)
	RoleEditor             // 1
	RoleAdmin
)

var roleNames = map[Role]string{
	RoleAdmin:  "admin",
	RoleEditor: "editor",
	RoleViewer: "viewer",
}

func (r Role) String() string {
	return roleNames[r]
}

func ParseRole(roleStr string) (Role, error) {
	for role, name := range roleNames {
		if name == roleStr {
			return role, nil
		}
	}
	return Role(0), fmt.Errorf("invalid role: %s", roleStr)
}

func (r *Role) Scan(value any) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case string:
		role, err := ParseRole(v)
		if err != nil {
			return err
		}
		*r = role
	case []byte:
		role, err := ParseRole(string(v))
		if err != nil {
			return err
		}
		*r = role
	default:
		return fmt.Errorf("cannot scan %T into Role", value)
	}
	return nil
}

func (r Role) Value() (driver.Value, error) { return r.String(), nil }

func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Role) UnmarshalJSON(data []byte) error {
	var roleStr string
	if err := json.Unmarshal(data, &roleStr); err != nil { // strips the quotes
		return err
	}
	role, err := ParseRole(roleStr)
	if err != nil {
		return err
	}
	*r = role
	return nil
}

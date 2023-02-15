package db

import "encoding/json"

func (user *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(user)
}

func (user *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, user)
}

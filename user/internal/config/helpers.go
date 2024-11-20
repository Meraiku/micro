package config

func ConvertToType(t string) RepoType {
	switch t {
	case "memory":
		return Memory
	}

	return Memory
}

package transform

func EnumMapKeys(enumMap map[string]int32, ignores ...int32) []string {
	out := make([]string, 0, len(enumMap))

enum:
	for k, v := range enumMap {
		for _, ign := range ignores {
			if v == ign {
				continue enum
			}
		}

		out = append(out, k)
	}

	return out
}

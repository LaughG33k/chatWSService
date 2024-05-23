package pkg

func CheckForAllKeys[V interface{}](data map[string]V, keys ...string) bool {

	for _, v := range keys {

		if _, ok := data[v]; !ok {
			return false
		}

	}

	return true
}

package openpipeline

func RemoveNils(m map[string]interface{}) {
	for a := range m {
		if m[a] == nil {
			delete(m, a)
		}
	}
}

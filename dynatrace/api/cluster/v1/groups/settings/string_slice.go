package groups

type strSlice []string

func (me strSlice) Nil() strSlice {
	if me == nil {
		return nil
	}
	if len(me) == 0 {
		return nil
	}
	return nil
}

type strSliceMap map[string][]string

func (me strSliceMap) Nil() strSliceMap {
	if me == nil {
		return nil
	}
	if len(me) == 0 {
		return nil
	}
	return nil
}

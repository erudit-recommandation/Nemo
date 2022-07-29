package test

func SliceContainTheSameElements[N any](s1 []N, s2 []N) bool {
	if len(s1) != len(s2) {
		return false
	}
	ss1 := map[interface{}]interface{}{}
	ss2 := map[interface{}]interface{}{}
	for _, el := range s1 {
		ss1[el] = el
	}

	for _, el := range s2 {
		ss2[el] = el
	}

	for _, v := range ss1 {
		if _, ok := ss2[v]; !ok {
			return false
		}
	}

	return true

}

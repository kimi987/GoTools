package atomic

func U46AddLimit(v *Uint64, toAdd, limit uint64) (uint64, uint64) {
	for {
		cur := v.Load()
		if cur >= limit {
			return 0, cur
		}

		maxCanAdd := limit - cur
		realAdd := toAdd
		if toAdd > maxCanAdd {
			realAdd = maxCanAdd
		}

		newAmount := cur + realAdd
		if v.CAS(cur, newAmount) {
			return realAdd, newAmount
		}
	}
}

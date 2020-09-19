package sharedguilddata

func GetLast(array []*Guild) *Guild {
	if len(array) > 0 {
		return array[len(array)-1]
	}

	return nil
}

func RemoveAndLeftShift(array []*Guild, removeId int64) []*Guild {
	if len(array) > 0 {
		for i, v := range array {
			if v.id == removeId {
				return leftShift(array, i)
			}
		}
	}

	return array
}

func leftShift(array []*Guild, startIndex int) []*Guild {

	n := len(array)
	for i := startIndex + 1; i < n; i++ {
		array[i-1] = array[i]
	}

	array[n-1] = nil
	return array[:n-1]
}

// id rank
type idRankSlice []*Guild

func (a idRankSlice) Len() int           { return len(a) }
func (a idRankSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a idRankSlice) Less(i, j int) bool { return a[i].id < a[j].id }

// prestige rank
type prestigeRankSlice []*Guild

func (a prestigeRankSlice) Len() int      { return len(a) }
func (a prestigeRankSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a prestigeRankSlice) Less(i, j int) bool {
	ai, aj := a[i], a[j]
	x, y := ai.YesterdayPrestige(), aj.YesterdayPrestige()
	if x != y {
		return x > y
	}
	return ai.id < aj.id
}

// prestigeCore rank
type prestigeCoreRankSlice []*Guild

func (a prestigeCoreRankSlice) Len() int      { return len(a) }
func (a prestigeCoreRankSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a prestigeCoreRankSlice) Less(i, j int) bool {
	ai, aj := a[i], a[j]
	x, y := ai.GetPrestigeCore(), aj.GetPrestigeCore()
	if x != y {
		return x > y
	}
	return ai.id < aj.id
}
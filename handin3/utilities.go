package handin3

type Vector struct {
	Clock []int32
}

func Max(a int32, b int32) int32 {
	if a < b {
		return b
	} else {
		return a
	}
}

func AdjustToOtherClock(personalClock Vector, otherClock Vector) Vector {
	ownLen := len(personalClock.Clock)
	otherLen := len(otherClock.Clock)
	sameLen := ownLen == otherLen
	if sameLen {
		for i := 0; i < otherLen; i++ {
			personalClock.Clock[i] = Max(personalClock.Clock[i], otherClock.Clock[i])
		}
	} else {
		dif := otherLen - ownLen

		for i := 0; i < dif; i++ {
			personalClock.Clock = append(personalClock.Clock, 0)
		}
		for i := 0; i < otherLen; i++ {
			personalClock.Clock[i] = Max(personalClock.Clock[i], otherClock.Clock[i])
		}
	}
	return personalClock
}

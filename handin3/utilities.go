package handin3

type Vector struct {
	clock []int
}

func Max(a int, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func UpdateClock(personalClock Vector, recievedClock Vector) Vector {
	ownLen := len(personalClock.clock)
	otherLen := len(recievedClock.clock)
	sameLen := ownLen == otherLen
	if sameLen {
		for i := 0; i < otherLen; i++ {
			personalClock.clock[i] = Max(personalClock.clock[i], recievedClock.clock[i])
		}
	} else {
		dif := otherLen - ownLen

		for i := 0; i < dif; i++ {
			personalClock.clock = append(personalClock.clock, 0)
		}
		for i := 0; i < otherLen; i++ {
			personalClock.clock[i] = Max(personalClock.clock[i], recievedClock.clock[i])
		}
	}
	return personalClock
}

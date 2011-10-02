package combatlog

type WindowFunc func(lastStart, start, lastEnd, end int)

func (cl CombatLog) RollingWindow(seconds, stepBy int, fun WindowFunc) {
	if len(cl) == 0 {
		return
	}

	var lastStart, start, lastEnd, end int

	threshold := int64(seconds) * 1e9
	startTime := cl[0].Time.Nanoseconds()

	for start < len(cl) {
		// Find the start of the window
		for _, e := range cl[start:] {
			time := e.Time.Nanoseconds()
			if time >= startTime {
				break
			}
			start++
		}
		// Find the end of the window
		if end < start {
			end = start
		}
		for _, e := range cl[end:] {
			time := e.Time.Nanoseconds()
			if time - startTime >= threshold {
				break
			}
			end++
		}
		fun(lastStart, start, lastEnd, end)
		lastStart, lastEnd = start, end
		startTime += int64(stepBy) * 1e9
	}
}

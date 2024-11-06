package litxaputil

import "slices"

func infix(p int, s ...string) Infix {
	return Infix{Pos: p, SyllableSplit: s}
}

// An Infix has data and behavior in order to insert itself.
type Infix struct {
	// Pos is for where it'll fit in.
	Pos int
	// SyllableSplit describes how the infix will be added.
	// first: will end the current syllable
	// middle: always its own syllable
	// last: will start the following syllable
	SyllableSplit []string
}

func (i Infix) Equal(o Infix) bool {
	return i.Pos == o.Pos && slices.Equal(i.SyllableSplit, o.SyllableSplit)
}

func (i Infix) Apply(curr []string, si, pos int) ([]string, int, int) {
	var tempData [32]string
	temp := tempData[:0]

	temp = append(temp, curr[:si+1]...)
	if len(i.SyllableSplit) > 2 {
		temp = append(temp, i.SyllableSplit[1:len(i.SyllableSplit)-1]...)
	}
	si2 := len(temp)
	temp = append(temp, "")
	temp = append(temp, curr[si+1:]...)

	temp[si] = curr[si][:pos] + i.SyllableSplit[0]
	temp[si2] = i.SyllableSplit[len(i.SyllableSplit)-1] + curr[si][pos:]

	curr = append(curr[:0], temp...)
	return curr, si + len(i.SyllableSplit) - 1, len(i.SyllableSplit[len(i.SyllableSplit)-1])
}

func ApplyInfixes(curr []string, infixNames []string, start int, stress int, positions [2][2]int) ([]string, int) {
	var infixes [3]*Infix
	for _, infixName := range infixNames {
		infix := FindInfix(infixName)
		if infix != nil {
			if infixes[infix.Pos] != nil && infix.Pos == 0 {
				if infixes[infix.Pos].Equal(infixMap["äp"]) && infix.Equal(infixMap["eyk"]) {
					infix = FindInfix("äpeyk")
				}
			}

			infixes[infix.Pos] = infix
		}
	}

	hasStressShift := stress == start && positions[0] == [2]int{0, 0}
	allInfixesTogether := positions[1] == positions[0]

	if infixes[0] != nil {
		next, si2, pos2 := infixes[0].Apply(curr, positions[0][0], positions[0][1])
		if !hasStressShift && stress == positions[0][0] {
			stress += si2 - positions[0][0]
		}

		positions = pushInfixPositions(positions, si2, pos2)
		curr = next
	}

	if infixes[1] != nil {
		next, si2, pos2 := infixes[1].Apply(curr, positions[0][0], positions[0][1])
		if !hasStressShift && stress >= positions[0][0] {
			stress += si2 - positions[0][0]
		}

		positions = pushInfixPositions(positions, si2, pos2)
		curr = next
	}

	if infixes[2] != nil {
		next, si2, _ := infixes[2].Apply(curr, positions[1][0], positions[1][1])
		if !hasStressShift && stress >= positions[1][0] {
			stress += si2 - positions[1][0]
		}

		curr = next
	}

	// Apply stress shift.
	if hasStressShift {
		if infixes[2] != nil && allInfixesTogether {
			stress += len(infixes[2].SyllableSplit) - 2
			if infixes[1] != nil {
				stress += len(infixes[1].SyllableSplit) - 1
			}
			if infixes[0] != nil {
				stress += len(infixes[0].SyllableSplit) - 1
			}
		} else if infixes[1] != nil {
			stress += len(infixes[1].SyllableSplit) - 2
			if infixes[0] != nil {
				stress += len(infixes[0].SyllableSplit) - 1
			}
		} else if infixes[0] != nil {
			stress += len(infixes[0].SyllableSplit) - 2
		}
	}

	for i, syllable := range curr {
		if i == 0 {
			continue
		}

		if syllable == "lll" || syllable == "rrr" {
			curr[i-1] += syllable[:1]
			curr = append(curr[:i], curr[i+1:]...)

			if stress >= i {
				stress -= 1
			}

			break
		}
	}

	return curr, stress
}

func pushInfixPositions(positions [2][2]int, si2 int, pos2 int) [2][2]int {
	if positions[0] == positions[1] {
		positions[1] = [2]int{si2, pos2}
	} else {
		positions[1] = [2]int{(si2 - positions[0][0]) + positions[1][0], positions[1][1]}
	}
	positions[0] = [2]int{si2, pos2}

	return positions
}

func FindInfix(name string) *Infix {
	infix, ok := infixMap[name]
	if !ok {
		return nil
	}

	return &infix
}

var infixMap = map[string]Infix{
	"äp":    infix(0, "ä", "p"),
	"eyk":   infix(0, "ey", "k"),
	"äpeyk": infix(0, "ä", "pey", "k"),

	"us": infix(1, "u", "s"),

	"am": infix(1, "a", "m"),
	"ìm": infix(1, "ì", "m"),
	"ìy": infix(1, "ì", "y"),
	"ay": infix(1, "a", "y"),

	"ìsy": infix(1, "ì", "sy"),
	"asy": infix(1, "a", "sy"),

	"er":  infix(1, "e", "r"),
	"arm": infix(1, "ar", "m"),
	"ìrm": infix(1, "ìr", "m"),
	"ìry": infix(1, "ìr", "y"),
	"ary": infix(1, "ar", "y"),

	"ol":  infix(1, "o", "l"),
	"alm": infix(1, "al", "m"),
	"ìlm": infix(1, "ìl", "m"),
	"ìly": infix(1, "ìl", "y"),
	"aly": infix(1, "al", "y"),

	"iv":   infix(1, "i", "v"),
	"irv":  infix(1, "ir", "v"),
	"ilv":  infix(1, "il", "v"),
	"imv":  infix(1, "im", "v"),
	"iyev": infix(1, "i", "ye", "v"),
	"ìyev": infix(1, "ì", "ye", "v"),

	"ei":  infix(2, "e", "i", ""),
	"eiy": infix(2, "e", "i", "y"),
	"eng": infix(2, "e", "ng"),
	"äng": infix(2, "ä", "ng"),
	"ats": infix(2, "a", "ts"),
	"uy":  infix(2, "u", "y"),
	"y":   infix(2, "", "y"), // In case the dict reports "uye'" as "u<y>e'"
}

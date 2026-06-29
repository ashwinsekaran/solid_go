package leetcode

// LeetCode 205 - Isomorphic Strings
//
// Pseudo code:
//   sToT = map, tToS = map
//   for each pair (cs, ct) at same index:
//     if sToT[cs] exists and sToT[cs] != ct: return false
//     if tToS[ct] exists and tToS[ct] != cs: return false
//     sToT[cs] = ct
//     tToS[ct] = cs
//   return true

func isIsomorphic(s string, t string) bool {
	sToT := make(map[byte]byte)
	tToS := make(map[byte]byte)
	for i := 0; i < len(s); i++ {
		cs, ct := s[i], t[i]
		if mapped, ok := sToT[cs]; ok && mapped != ct {
			return false
		}
		if mapped, ok := tToS[ct]; ok && mapped != cs {
			return false
		}
		sToT[cs] = ct
		tToS[ct] = cs
	}
	return true
}

package util

//verifica se a primeira string contém a segunda, se sim, retorna true
// e a posição da segunda string dentro da primeira, caso contrário, retorna false e -1
func Contains(slice string, item string) (bool, int) {
	pos := -1
	for i := 0; i < len(slice); i++ {
		if slice[i] == item[0] {
			match := true
			for j := 1; j < len(item); j++ {
				if i+j >= len(slice) || slice[i+j] != item[j] {
					match = false
					break
				}
			}
			if match {
				pos = i
				return true, pos
			}
		}
	}
	return false, pos

}

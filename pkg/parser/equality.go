package parser

import (
	"hash/fnv"
)

func hashAll(vars []string) map[string]int {
	t := make(map[string]int)
	for _, v := range vars {
		t[v] = int(hash(v))
	}
	return t
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// IsTrue returns whether the equation is true or not linear equations only
// func IsTrue(equation string) (bool, error) {
// 	parts := strings.Split(equation, "=")

// 	if len(parts) < 2 {
// 		return false, fmt.Errorf("Invalid equation, missing '='")
// 	}

// 	trees := []*AST{}
// 	vals1 := []int{}
// 	vals2 := []int{}
// 	for _, exp := range parts {
// 		tree, err := Parse(exp)
// 		if err != nil {
// 			return false, err
// 		}
// 		trees = append(trees, tree)
// 		ev1 := tree.Evaluate(hashAll(tree.Symbols()))
// 		if ev1 == nil || ev1.Val == nil {
// 			return false, fmt.Errorf("Unknown error")
// 		}
// 		ev2 := tree.Evaluate(hashAll(tree.Symbols()))
// 		if ev2 == nil || ev2.Val == nil {
// 			return false, fmt.Errorf("Unknown error")
// 		}

// 		vals1 = append(vals1, int(*ev1.Val))
// 		vals2 = append(vals2, int(*ev2.Val))
// 	}

// 	for i := 0; i < len(vals1)-1; i++ {
// 		if vals1[i] != vals1[i+1] || vals2[i] != vals2[i+1] {
// 			return false, nil
// 		}
// 	}
// 	return true, nil
// }

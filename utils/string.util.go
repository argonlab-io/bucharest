package utils

import "sort"

func ArrayStringContians(s []string, searchterm string) (int, bool) {
    i := sort.SearchStrings(s, searchterm)
    return i, i < len(s) && s[i] == searchterm
}


package utils

func IsArrayIncludes(array []string, element string) int {
    for position, value := range array {
        if value == element {
            return position
        }
    }
    return -1
}

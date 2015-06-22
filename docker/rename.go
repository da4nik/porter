package docker

func Rename(oldName, newName string) {
    logger.Printf("Rename container %s to %s (code: %d)\n", oldName, newName, renameContainer(oldName, newName))
}

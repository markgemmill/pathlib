package pathlib

import (
	"os"
)

var (
	PRIVATE_FILE = os.FileMode(0600)
	PRIVATE_EXE  = os.FileMode(0700)
	PRIVATE_DIR  = os.FileMode(0700)

	READONLY_FILE = os.FileMode(0644)
	READONLY_DIR  = os.FileMode(0755)

	PUBLIC_FILE = os.FileMode(0666)
	PUBLIC_DIR  = os.FileMode(0777)
)

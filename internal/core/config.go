package core

import (
	"fmt"

	"github.com/ramzeng/ai-endpoint/package/config"
)

func ReadConfig(path string) {
	if err := config.Initialize(path); err != nil {
		fmt.Print("[Core]: failed to read the configuration file: ", err)
	}
}

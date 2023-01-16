package monitoring

import (
	"fmt"
	"log"
	"os"
)

func Init() {
	log.SetPrefix(fmt.Sprintf("[%d]%s: ", os.Getpid(), "Roberts Concordance"))
}

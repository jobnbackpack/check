package api

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

func WriteToFile(goals []textinput.Model) {
	var x = []byte{}

	for i := 0; i < len(goals); i++ {
		b := []byte("Goal " + fmt.Sprint(i+1) + " " + goals[i].Value() + "\n")
		for j := 0; j < len(b); j++ {
			x = append(x, b[j])
		}
	}

	os.WriteFile(time.Now().Format("2006-01-02")+".txt", x, 0644)
}

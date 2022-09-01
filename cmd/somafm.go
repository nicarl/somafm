package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/nicarl/somafm/audio"
	"github.com/nicarl/somafm/prompt"
	"github.com/nicarl/somafm/radioChannels"
)

func main() {
	//playMusic()
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	fmt.Println(text)
}

func playMusic() {
	radioCh, err := radioChannels.GetChannels()
	if err != nil {
		log.Fatal(err)
	}
	selectedCh, err := prompt.SelectChannel(radioCh)
	if err != nil {
		log.Fatal(err)
	}
	streamUrl, err := radioChannels.GetStreamUrl(selectedCh)
	if err != nil {
		log.Fatal(streamUrl)
	}
	err = audio.PlayRemoteFile(streamUrl)
	if err != nil {
		log.Fatal(err)
	}
}

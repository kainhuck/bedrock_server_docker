package main

import (
	"log"
	"os"
	"path"
	"text/template"

	"github.com/gocolly/colly"
)

const (
	BedrockDownloadPage        = "https://www.minecraft.net/en-us/download/server/bedrock"
	BedrockLinkSelector        = "a[data-platform='serverBedrockLinux']"
	BedrockLinkPreviewSelector = "a[data-platform='serverBedrockPreviewLinux']"
)

var (
	InstallRootDir = "/opt"
	WorkDir        = path.Join(os.TempDir(), "minecraft_kainhuck")
)

func main() {
	link, err := GetBedrockDownloadLink(BedrockLinkSelector)
	if err != nil {
		log.Fatalf("get link failed: %v", err)
	}

	if err := TemplateDockerfile(link); err != nil {
		log.Fatalf("get dockerfile failed: %v", err)
	}
	log.Println("SUCCESS")
}

func TemplateDockerfile(link string) error {

	dockerfileTemp, err := template.New("xxxx").Parse(DockerfileTemp)
	if err != nil {
		return err
	}

	f, err := os.OpenFile("Dockerfile", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	return dockerfileTemp.Execute(f, link)
}

func GetBedrockDownloadLink(selector string) (link string, err error) {
	log.Println("get bedrock download link ...")
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/109.0"))
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		link = e.Attr("href")
	})

	err = c.Visit(BedrockDownloadPage)
	return
}

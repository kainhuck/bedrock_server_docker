package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"text/template"

	"github.com/gocolly/colly"
)

const (
	BedrockDownloadPage        = "https://www.minecraft.net/en-us/download/server/bedrock"
	BedrockLinkSelector        = "a[data-platform='serverBedrockLinux']"
	BedrockLinkPreviewSelector = "a[data-platform='serverBedrockPreviewLinux']"
	DefaultBedrockLink         = "https://minecraft.azureedge.net/bin-linux/bedrock-server-1.19.52.01.zip"
)

var (
	InstallRootDir = "/opt"
	WorkDir        = path.Join(os.TempDir(), "minecraft_kainhuck")
	Version        = ""
	DownloadLink   = "https://minecraft.azureedge.net/bin-linux/bedrock-server-%s.zip"
)

func init() {
	flag.StringVar(&InstallRootDir, "i", InstallRootDir, "-i <安装目录>")
	flag.StringVar(&Version, "v", "", "-v <安装版本>")
	flag.Parse()
}

func main() {
	// 0. 检查环境docker，docker-compose是否安装
	CheckEnv()
	// 1. 进入工作目录，打包镜像
	defer NewDir(WorkDir)()
	// 	1.1 如果用户指定了版本则使用用户指定的版本，否则去网站上拉去，如果拉取失败则使用默认的url
	if Version != "" {
		DownloadLink = fmt.Sprintf(DownloadLink, Version)
	} else {
		link, _ := GetBedrockDownloadLink(BedrockLinkSelector)
		if link == "" {
			DownloadLink = DefaultBedrockLink
		} else {
			DownloadLink = link
		}
	}
	// 	1.2 新建Dockerfile，构建docker镜像
	if err := TemplateDockerfile(DownloadLink); err != nil {
		log.Fatal(err)
	}
	version, err := getVersion(DownloadLink)
	if err != nil {
		log.Fatal(err)
	}
	image := fmt.Sprintf("kainhuck/bedrock:%s", version)
	runCmd(fmt.Sprintf("docker build -t %s -f %s %s", image, filepath.Join(WorkDir, "Dockerfile"), WorkDir))
	// 	1.3 删除工作目录

	// 2. 创建安装目录
	installDir := filepath.Join(InstallRootDir, "minecraft_kainchuk")
	NewDir(installDir)
	// 	2.1 生成各种配置文件
	NewDir(filepath.Join(installDir, "worlds"))
	if err := TemplateBedrock(installDir, image); err != nil {
		log.Fatal(err)
	}
	//	2.2 启动服务

	// link, err := GetBedrockDownloadLink(BedrockLinkSelector)
	// if err != nil {
	// 	log.Fatalf("get link failed: %v", err)
	// }

	// if err := TemplateDockerfile(link); err != nil {
	// 	log.Fatalf("get dockerfile failed: %v", err)
	// }
	// log.Println("SUCCESS")
}

func TemplateDockerfile(link string) error {

	dockerfileTemp, err := template.New("docker").Parse(DockerfileTemp)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join(WorkDir, "Dockerfile"), os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	return dockerfileTemp.Execute(f, link)
}

func TemplateBedrock(installDir string, image string) error {

	dockerCompose := DockerCompose{
		Image:      image,
		InstallDir: installDir,
	}

	dockercomposeTemp, err := template.New("docker-compose").Parse(DockercomposeTemp)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join(installDir, "docker-compose.yml"), os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	if err := dockercomposeTemp.Execute(f, dockerCompose); err != nil {
		return err
	}

	return nil
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

func CheckEnv() {
	if _, err := exec.LookPath("docker"); err != nil {
		log.Fatal("docker not install")
	}
	if _, err := exec.LookPath("docker-compose"); err != nil {
		log.Fatal("docker-compose not install")
	}
}

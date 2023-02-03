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
	Xuid           = ""
	Mode           = "survival"
	Difficulty     = "hard"
	WorldName      = "new_world"
	WorldSeed      = ""
)

func init() {
	flag.StringVar(&InstallRootDir, "i", InstallRootDir, "-i <安装目录>")
	flag.StringVar(&Version, "v", "", "-v <安装版本>")
	flag.StringVar(&Xuid, "x", "", "-x <xuid>")
	flag.StringVar(&Mode, "m", Mode, "-m <模式 survival creative adventure>")
	flag.StringVar(&Difficulty, "d", Difficulty, "-d <难度 peaceful, easy, normal, hard>")
	flag.StringVar(&WorldName, "n", WorldName, "-n <世界名称>")
	flag.StringVar(&WorldSeed, "s", WorldSeed, "-s <种子>")
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
	if err := RunCmd(fmt.Sprintf("docker build -t %s -f %s %s", image, filepath.Join(WorkDir, "Dockerfile"), WorkDir)); err != nil {
		log.Fatal(err)
	}
	// 	1.3 删除工作目录

	// 2. 创建安装目录
	installDir := filepath.Join(InstallRootDir, "minecraft_kainchuk")
	NewDir(installDir)
	// 	2.1 生成各种配置文件
	NewDir(filepath.Join(installDir, "worlds"))
	if err := TemplateBedrock(installDir, image); err != nil {
		log.Fatal(err)
	}

	Hello(image, installDir, version)
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

	temp := func(filename string, tpl string, field interface{}) error {
		t, err := template.New(filename).Parse(tpl)
		if err != nil {
			return err
		}

		f, err := os.OpenFile(filepath.Join(installDir, filename), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		return t.Execute(f, field)
	}

	if err := temp("docker-compose.yml", DockercomposeTemp, DockerCompose{
		Image:      image,
		InstallDir: installDir,
	}); err != nil {
		return err
	}

	if err := temp("permissions.json", PermissionsJsonTemp, PermissionsJson{
		XUID: Xuid,
	}); err != nil {
		return err
	}

	if err := temp("server.properties", ServerPropertiesTemp, ServerProperties{
		Mode:       Mode,
		Difficulty: Difficulty,
		WorldName:  WorldName,
		WorldSeed:  WorldSeed,
	}); err != nil {
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

func Hello(image string, installDir string, version string) {
	fmt.Println("=============================")
	fmt.Println("恭喜，服务部署成功！")
	fmt.Printf("镜像名称: %s\n", image)
	fmt.Printf("安装路径: %s\n", installDir)
	fmt.Printf("版本: %s\n", version)
	fmt.Printf("世界名称: %s\n", WorldName)
	fmt.Printf("世界模式: %s\n", Mode)
	fmt.Printf("世界难度: %s\n", Difficulty)
	if len(Xuid) > 0 {
		fmt.Printf("世界管理员: %s\n", Xuid)
	}
	if len(WorldSeed) > 0 {
		fmt.Printf("世界种子: %s\n", WorldSeed)
	}
	fmt.Println("=============================")
	fmt.Printf("运维指南:\n")
	fmt.Printf("启动服务: docker-compose -f %s up -d\n", filepath.Join(installDir, "docker-compose.yml"))
	fmt.Printf("服务暂停: docker-compose -f %s stop\n", filepath.Join(installDir, "docker-compose.yml"))
	fmt.Printf("服务删除: docker-compose -f %s rm\n", filepath.Join(installDir, "docker-compose.yml"))
	fmt.Printf("服务状态: docker-compose -f %s ps\n", filepath.Join(installDir, "docker-compose.yml"))
	fmt.Println("=============================")
	fmt.Printf("卸载指南:\n")
	fmt.Printf("1. 执行: docker-compose -f %s stop\n", filepath.Join(installDir, "docker-compose.yml"))
	fmt.Printf("2. 执行: docker-compose -f %s rm -f\n", filepath.Join(installDir, "docker-compose.yml"))
	fmt.Printf("3. 执行: sudo rm -rf %s\n", installDir)
	fmt.Printf("4. 执行: docker image rm %s\n", image)
}

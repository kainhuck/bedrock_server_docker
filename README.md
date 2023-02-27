# docker image for bedrock_server

自动部署我的世界基岩版，利用docker-compose管理，简单方便

## Usage

```bash
Usage of bsd:
  -d string
        -d <难度 peaceful, easy, normal, hard> (default "hard")
  -i string
        -i <安装目录> (default "/opt")
  -m string
        -m <模式 survival creative adventure> (default "survival")
  -n string
        -n <世界名称> (default "new_world")
  -s string
        -s <种子>
  -u    更新 -v 指定的版本
  -v string
        -v <安装版本>
  -x string
        -x <xuid>
```

## 工作目录结构
```perl
minecraft_kainchuk
├── docker-compose.yml // docker-compose
├── permissions.json   // 权限文件
├── server.properties  // 配置文件
└── worlds             // 世界目录
```

enjoy it 😊
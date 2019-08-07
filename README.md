# Project setup

## 创建 Personal Access Token 并登录

```shell
docker login registry.gitlab.botpy.com
```

## 拉取镜像

```shell
docker pull registry.gitlab.botpy.com/pengfei/yobee-debugtool
```

## 运行

```shell
docker run -v PATH/TO/CONFIG_SYSTEM:/server/configs/systems.toml -p 5000:5000 --name yobee-debugtool -d registry.gitlab.botpy.com/pengfei/yobee-debugtool
```

# 使用帮助

## 保存文件

Bot 接受两种消息: 文件和链接.

支持以下链接:

1. 公开频道 (具有用户名) 的消息链接, 例如: `https://t.me/acherkrau/1097`. **即使频道禁止了转发和保存, Bot 依然可以下载其文件.**
2. Telegra.ph 的文章链接, Bot 将下载其中的所有图片

## 静默模式 (silent)

使用 `/silent` 命令可以开关静默模式.

默认情况下不开启静默模式, Bot 会询问你每个文件的保存位置.

开启静默模式后, Bot 会直接保存文件到默认位置, 无需确认.

在开启静默模式之前, 需要使用 `/storage` 命令设置默认保存位置.

## Stream 模式

在配置文件中将 `stream` 设置为 `true` 可以开启 Stream 模式.

未开启时, Bot 处理任务分为两步: 下载和上传. Bot 会将文件暂存到本地, 然后上传到对应存储位置, 最后删除本地文件.

开启后, Bot 将直接将文件流式传输到存储端, 不需要下载到本地.

该功能对于硬盘空间有限的部署环境十分有用, 然而相较于普通模式也具有一些弊端:

- 无法使用多线程从 telegram 下载文件, 速度较慢.
- 网络不稳定时, 任务失败率高.
- 无法在中间层对文件进行处理, 例如自动文件类型识别.

**不支持** Stream 模式的存储端:

- alist
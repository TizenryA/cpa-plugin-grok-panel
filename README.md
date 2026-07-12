# Grok Panel for CLIProxyAPI

Grok Panel 是 [`router-for-me/CLIProxyAPI`](https://github.com/router-for-me/CLIProxyAPI) 的原生插件，用于统计和管理当前 CPA 实例中的 Grok / xAI auth 文件。

> 当前版本：`v1.1.0`｜目标平台：Linux amd64｜界面语言：中文

## 用途与功能

- Grok 文件、活跃/禁用账号、成功/失败请求统计；
- 按“成功请求 × 平均 Token”估算使用量与剩余容量；
- 请求趋势、搜索、筛选、排序和账号明细；
- 账号分类：`Free`、`Super`、`Heavy`、`Unknown`；
- 手动检查单个、选中或当前可见账号；
- 健康状态：健康、禁用、明确失效、暂时异常、未知；
- 单个删除、批量删除、清理已确认失效账号；
- 自动检查和自动删除开关默认关闭；
- Super、Heavy、Unknown 默认禁止自动删除；
- 手动删除使用二次点击，不使用弹窗；
- 429、5xx、超时不会被当作失效账号；
- 只有连续 3 次明确 401/403 才进入自动删除候选。

Token 为估算值，不是 xAI 返回的真实计费数据。默认每账号容量 `2,000,000 Token`、每次成功请求 `5,000 Token`，可在面板内调整。

## 安装方式 A：插件商店（推荐）

在 CPA 的 `config.yaml` 中启用插件并添加本仓库源：

```yaml
plugins:
  enabled: true
  dir: "plugins"
  store-sources:
    - "https://raw.githubusercontent.com/TizenryA/cpa-plugin-grok-panel/main/registry.json"
  configs:
    grok-panel:
      enabled: true
```

重启 CPA，进入管理中心的插件商店，安装 **Grok Panel**。安装完成后，在插件菜单打开 **Grok Panel**。

也可以使用 CPA 管理 API 安装：

```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_MANAGEMENT_KEY" \
  "https://YOUR_CPA_HOST/v0/management/plugin-store/grok-panel/install?version=v1.1.0"
```

这里的管理密钥只用于 CPA 执行安装，不会写入插件。

## 安装方式 B：手动安装

从 [Releases](https://github.com/TizenryA/cpa-plugin-grok-panel/releases) 下载与宿主匹配的压缩包，例如：

```text
grok-panel_1.1.0_linux_amd64.zip
```

解压后将 `grok-panel.so` 放入 CPA 配置的插件目录：

```text
plugins/grok-panel.so
```

确保配置中已启用插件，然后重启 CPA。

## 不需要重复填写管理密钥

插件读取统计数据时使用 CPA 官方 host callback：

```text
host.auth.list
host.auth.get
host.auth.get_runtime
```

因此不需要填写 CPA 地址或密钥。实际删除时，由于 CPA 暂未提供 `host.auth.delete`，面板会复用 CPA 管理中心当前保存的连接凭据，调用官方接口：

```text
DELETE /v0/management/auth-files
```

安全边界：

- 源码、Release、插件配置均不包含任何管理密钥；
- 密钥只在当前浏览器中用于请求当前 CPA；
- 删除仍受 CPA 官方管理鉴权保护；
- 如果登录时没有勾选“记住密码”，检查和删除会提示重新登录并保存会话；
- 不向浏览器返回 access token、refresh token、SSO token 或 cookie。

## 账号分类

插件从 CPA auth 元数据和 auth JSON 的套餐字段中识别：

```text
Free     普通账号
Super    SuperGrok / 高级套餐
Heavy    Heavy 套餐
Unknown  无可靠套餐信号
```

无法可靠判断时必须归入 `Unknown`，不会猜成普通账号。Unknown 默认受自动删除保护。

## 健康检查与删除规则

健康检查优先使用 CPA runtime 状态和明确的认证错误：

- 明确 `401 Unauthorized` 或 `403 Forbidden`：累计失效次数；
- `429`、`5xx`、网络超时：只视为暂时异常，不累计失效；
- 连续 3 次明确 401/403：成为自动删除候选；
- 任意一次健康结果：失效连续次数归零；
- Super、Heavy、Unknown：默认永不自动删除；
- 自动检查、自动删除：均默认关闭，由用户主动开启。

手动删除不要求账号先达到失效阈值，但仍遵守当前保护开关，并要求二次点击确认。

## 设置

面板设置保存在当前浏览器：

- 每账号估算容量；
- 平均 Token/请求；
- 连续认证失败阈值；
- 自动检查；
- 自动删除；
- Super / Heavy / Unknown 保护。

性能相关功能默认关闭，不会安装后立即批量检查账号。

## 升级与卸载

升级：在 CPA 插件商店点击更新，或重新调用安装接口指定新版本。

卸载：在 CPA 插件管理页面卸载 `grok-panel`，然后重启或 reload 插件系统。

## 构建

需要与 CPA 宿主兼容的 Linux、CPU 架构及 Go 工具链：

```bash
go test ./...
go vet ./...
go build -buildmode=c-shared -o grok-panel.so .
```

## 隐私与安全

- 不连接作者的 CPA；
- 每个安装者只读取自己的 CPA 数据；
- 不上传账号文件或凭据；
- 不在日志中打印 token、cookie 或授权头；
- 删除使用 CPA 官方管理 API，不直接遍历或修改宿主文件系统。

## License

MIT

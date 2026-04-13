# Toybox 使用示例

## 示例列表

### 1. basic — 基础用法

最简单的用法：自定义一个 HTTP Transport，交给 ToyBox 管理生命周期（启动、信号处理、优雅关闭）。

```bash
cd basic && go run .
# 访问 http://localhost:8080/hello
```

### 2. config — 配置管理与热重载

演示从本地 YAML 文件加载配置，Transport 通过实现 `ConfigureMatcher` 接口自动接收配置注入，并支持文件变更时自动热重载。

```bash
cd config && go run .
# 修改 config.yaml 后服务会自动重载
```

### 3. multi_transport — 多 Transport 并发运行

同时运行多个 Transport（如 API 服务 + Metrics 服务），ToyBox 会并发启动它们并统一管理生命周期。

```bash
cd multi_transport && go run .
# API:     http://localhost:8080/api/users
# Metrics: http://localhost:9090/metrics
```

### 4. bootstrap — Bootstrap 一键初始化

使用 `bootstrap.Config` 声明式配置，自动初始化配置源和服务注册中心。支持从命令行参数、环境变量和默认值合并配置。

```bash
cd bootstrap && go run .
# 或通过环境变量:
# TOYBOX_SERVICE_NAME=my-svc go run .
```

## 核心概念

| 概念 | 说明 |
|------|------|
| `Transport` | 实现 `Start(ctx) error` 和 `Stop(ctx) error` 的服务抽象 |
| `ConfigureMatcher` | Transport 可选实现，通过 `TransportName()` 声明自己需要哪段配置 |
| `Registrar` | 服务注册/注销接口，内置 etcd、consul、zookeeper 实现 |
| `ConfigMate` | 配置管理接口，支持 YAML/JSON/TOML 格式 |
| `Bootstrap` | 声明式启动配置，一键关联注册中心和配置源 |

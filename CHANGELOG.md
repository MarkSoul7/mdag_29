# 更新日志

本项目的所有重要变更都记录在此文件中，格式参考
[Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### 文档

- 校准代码量（核心约 600 行，含测试约 900 行）、README 行数、docs 篇数与五看三定字数
- 补全 README 文档索引与 1.0.0 文档清单，补充“示例 CID 已截断”的说明
- 分支管理文档按 `git merge` 实际输出补全冲突快照与嵌套代码围栏
- CONTRIBUTING 的合并策略与仓库保留分支的例外说明对齐

### 计划中

- `mdag fsck`：全量校验对象完整性
- `mdag gc`：清理不可达对象
- `add --chunk-size`：把 List 分块接入 CLI
- 只读 HTTP 网关

## [1.0.0] - 2026-11-09

### 新增

- Blob / Tree / List 三种 Merkle DAG 对象，CID = SHA-256(JSON(Object))
- 内存存储与文件存储两种 `Store` 实现（`./data/<cid>.json`）
- 文件与目录递归导入，目录项按名称排序保证确定性
- 大文件定长分块（List）与还原读取
- 路径解析、文件读取、目录列举（`resolve` / `cat` / `ls`）
- CLI 子命令 `add`、`resolve`、`cat`、`ls`，支持 `--mem`
- 10 个自动化测试，覆盖 CID 稳定性、变更传播、错误路径、分块读取
- GitHub Actions CI：gofmt、go vet、go test
- 完整文档：README 与 `docs/` 六篇（架构设计、五看三定分析、分支管理与冲突处理、协作记录、作业要求对照、Issue 讨论提纲）
- Issue / Pull Request 模板与贡献指南
- MIT 许可证

### 修复

- 修正 README 中失效的代码块围栏与示例命令
- 统一代码格式（gofmt）

[Unreleased]: https://github.com/MarkSoul7/mdag_29/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/MarkSoul7/mdag_29/releases/tag/v1.0.0

# 贡献指南

感谢你关注 mdag_29！这是一个面向教学与学习的小型开源项目，欢迎通过 Issue 和
Pull Request 参与改进。

## 一、行为准则

- 尊重不同经验水平的参与者，对事不对人；
- 提问前先搜索已有 Issue，避免重复；
- 提交代码前请自测，描述问题时给出复现步骤与环境信息。

## 二、开发环境

```bash
git clone https://github.com/MarkSoul7/mdag_29.git
cd mdag_29
go version        # 需要 1.24+
go test ./...     # 确认基线通过
```

本项目**不引入第三方依赖**。如需新增依赖，请先在 Issue 中说明理由。

## 三、分支规范

| 分支 | 用途 | 示例 |
|------|------|------|
| `main` | 稳定分支，始终可构建、测试通过 | — |
| `feature/<name>` | 新功能 | `feature/fsck` |
| `fix/<name>` | 缺陷修复 | `fix/resolve-empty-path` |
| `docs/<name>` | 文档修改 | `docs/quickstart` |
| `chore/<name>` | 构建、CI、杂项 | `chore/ci-go-version` |

请从最新的 `main` 切出分支，不要在 `main` 上直接提交。

## 四、提交信息规范

采用 [Conventional Commits](https://www.conventionalcommits.org/)：

```text
<type>(<scope>): <subject>

<body 可选>

<footer 可选>
```

| type | 含义 |
|------|------|
| `feat` | 新功能 |
| `fix` | 修复缺陷 |
| `docs` | 仅文档 |
| `test` | 仅测试 |
| `refactor` | 重构（不改变行为） |
| `chore` | 构建、CI、依赖、杂项 |

示例：

```text
feat(cli): 增加 add --chunk-size 参数
fix(resolver): 修正空路径解析返回根对象类型
docs(readme): 补充 Windows 下的路径转换说明
```

## 五、代码风格

- 提交前执行 `gofmt -w .`，CI 会检查 `gofmt -l .` 是否为空；
- 执行 `go vet ./...` 无告警；
- 导出标识符必须有注释，注释以标识符名开头；
- 错误使用 `%w` 包装，保留错误链；
- 不引入不必要的抽象，优先可读性。

## 六、测试要求

- 新增功能必须配套测试；
- 修复缺陷应先补一个能复现的失败测试；
- 本地运行 `go test -v -count=1 ./...` 全部通过后再提交；
- 不要为了让测试通过而删除断言或跳过用例（确需跳过请说明原因）。

## 七、Pull Request 流程

1. Fork 仓库（或在本仓库创建分支）；
2. 从 `main` 切出 `feature/`、`fix/` 或 `docs/` 分支；
3. 提交变更，确保 CI 通过；
4. 发起 PR，填写模板中的“变更内容 / 关联 Issue / 验证方式”；
5. 至少 1 位维护者评审通过后，用 `git merge --no-ff` 合入（保留合并节点）；单提交的文档微调可用 **Squash merge**；
6. 合并后可删除远端分支。本仓库为便于复核保留已合并分支，实际项目按团队规范清理。

### PR 检查清单

- [ ] `gofmt -l .` 输出为空
- [ ] `go vet ./...` 通过
- [ ] `go test -count=1 ./...` 全部通过
- [ ] 新增/修改的行为有测试或文档说明
- [ ] 提交信息符合规范
- [ ] 没有提交编译产物（`mdag.exe`）或 `data/` 目录

## 八、Issue 规范

- 使用对应的 Issue 模板（Bug / 功能建议 / 提问）；
- Bug 报告请附：操作系统、Go 版本、完整命令、实际输出、期望输出；
- 功能建议请说明使用场景，而不只是“希望增加 X”。

## 九、文档贡献

文档与代码同等重要。修改文档时请：

- 保持命令输出与当前代码一致（可实际运行验证）；
- 中英文术语首次出现时给出全称；
- 不要粘贴无法复现的“预期输出”。

## 十、许可证

提交贡献即表示你同意以 [MIT License](LICENSE) 授权你的贡献。

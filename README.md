# mdag_29 · 简化版 Merkle DAG 文件存储与路径解析系统

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green)
![Tests](https://img.shields.io/badge/tests-10%20passed-brightgreen)
![Course](https://img.shields.io/badge/%E8%AF%BE%E7%A8%8B-%E5%BC%80%E6%BA%90%E6%96%87%E5%8C%96%E8%BF%90%E5%8A%A8-blue)

> 《开源文化运动》期末考核项目  
> 区块链工程 231 · 杨琦航 · 2023131029  
> 仓库名 `mdag_29`（项目名称 + 学号后两位）

---

## 目录

- [一、项目简介](#一项目简介)
- [二、快速开始](#二快速开始)
- [三、命令与示例](#三命令与示例)
- [四、核心原理](#四核心原理)
- [五、项目结构](#五项目结构)
- [六、测试](#六测试)
- [七、文档索引](#七文档索引)
- [八、分支模型与协作](#八分支模型与协作)
- [九、路线图](#九路线图)
- [十、许可证](#十许可证)

---

## 一、项目简介

**mdag** 是一个用 Go 语言实现的**简化版内容寻址文件系统**（Content-Addressed
Storage）。它把文件编码为 **Blob**、把目录编码为 **Tree**、把可选的大文件分块编码为
**List**，并用对象序列化内容的 SHA-256 摘要作为内容标识（CID）；给定一个根 CID
和一个路径，就能逐层解析并读回目标文件。

整个项目**零第三方依赖**，只使用 Go 标准库，代码量约 700 行，配套 10 个自动化测试。

| 特性 | 说明 |
|------|------|
| 内容寻址 | CID = `hex(SHA256(JSON(Object)))`，内容相同则 CID 相同 |
| 不可变 | 对象一旦写入即不可修改，修改文件必然产生新 CID |
| 自动去重 | 相同内容只存一份，重复导入是幂等操作 |
| 可验证 | 任一节点被篡改，其 CID 及所有祖先 CID 都会改变 |
| 可解释 | 每一层数据结构都能对应到源码与文档，适合教学与答辩 |
| 可替换存储 | `Store` 接口抽象，内置内存存储与文件存储两种实现 |

### 为什么叫 Merkle DAG

Tree 对象通过 `Link` 指向子对象的 CID，形成“父哈希包含子哈希”的结构：子节点变化会
逐层向上传播，最终改变根 CID。由于每个对象只有一个 CID、且 Link 有方向，这个图天然
无环，因此是 **Merkle DAG**（有向无环图），而不是普通的 Merkle Tree。

---

## 二、快速开始

### 2.1 环境要求

- Go 1.24 或更高版本（`go version` 查看）
- Git
- 不需要任何第三方依赖，不需要数据库

### 2.2 获取与构建

```bash
git clone https://github.com/MarkSoul7/mdag_29.git
cd mdag_29
go build -o mdag ./cmd/mdag
```


#### Windows 说明

在 PowerShell 或 CMD 中直接使用上面的命令即可。若在 Git Bash 中执行 `resolve`、`cat`、
`ls` 等子命令，请先运行 `export MSYS_NO_PATHCONV=1`，否则 `/docs` 会被 MSYS 改写成
Windows 绝对路径，导致 "name C: not found in tree" 错误。
Windows 下建议把产物命名为 `mdag.exe`：

```bash
go build -o mdag.exe ./cmd/mdag
```

### 2.3 导入示例目录

```bash
./mdag add testdata/demo
```

输出（实测）：

```text
Root CID: ae033b4c1ea3a83944ec106b17e2d080c19848e5f110a930b8e60480083b45b5
```

### 2.4 按路径解析

```bash
./mdag resolve ae033b4c1ea3a83944ec106b17e2d080c19848e5f110a930b8e60480083b45b5 /docs/report.txt
```

输出：

```text
Target CID: 9f9b877e2a702f761b4bba5a78bfcf382820a7392babd402a2d685bf84173dae
Type: blob
```

### 2.5 读取与列目录

```bash
./mdag cat ae033b4c1ea3a83944ec106b17e2d080c19848e5f110a930b8e60480083b45b5 /docs/report.txt
./mdag ls  ae033b4c1ea3a83944ec106b17e2d080c19848e5f110a930b8e60480083b45b5 /
```

输出：

```text
Course project report content.

README.md  blob  cfbefcec90027020a6255f9dd6dd67c33c0e4db20eeed5159818448db27aaed4  64
data       tree  b65593877e2428a22863083d97e723bfeab675123f90d6745849e0ecd9a497c7  0
docs       tree  1a13f7573ff7ee2ea801c398e2c539e8bf61f9d3cd5522c006e0c92401a58872  0
```

> **Windows 提示**：在 Git Bash 中执行时，`/docs/report.txt` 会被 MSYS 误当成
> Windows 绝对路径改写。请先执行 `export MSYS_NO_PATHCONV=1`，或改用 PowerShell /
> CMD 运行。

---

## 三、命令与示例

### 3.1 命令总览

| 命令 | 作用 | 示例 |
|------|------|------|
| `mdag add <path>` | 导入文件或目录，输出根 CID | `mdag add testdata/demo` |
| `mdag resolve <root> <path>` | 解析路径，输出目标 CID 与类型 | `mdag resolve <root> /docs` |
| `mdag cat <root> <path>` | 输出文件内容 | `mdag cat <root> /README.md` |
| `mdag ls <root> <path>` | 列出目录子项（名称 / 类型 / CID / 大小） | `mdag ls <root> /docs` |
| `--mem` | 全局可选参数，使用内存存储而非 `./data` | `mdag add --mem testdata/demo` |

退出码：成功为 `0`，任何错误为 `1`，错误信息输出到 stderr。

### 3.2 列目录

```bash
./mdag ls ae033b4c1ea3a83944ec106b17e2d080c19848e5f110a930b8e60480083b45b5 /docs
```

```text
notes.txt   blob  04a66e51e4ed8a86fdaa976cf239f1697c4b686979294c5f8e8c9728d3413952  26
report.txt  blob  9f9b877e2a702f761b4bba5a78bfcf382820a7392babd402a2d685bf84173dae  31
```

### 3.3 错误处理

```bash
# 1. 路径不存在
./mdag resolve <root> /docs/missing.txt
# Error: path "/docs/missing.txt": name "missing.txt" not found in tree 1a13f757...

# 2. 把目录当文件读取
./mdag cat <root> /docs
# Error: "/docs" (1a13f757...) is a tree, not a file

# 3. CID 不存在
./mdag resolve deadbeef /x
# Error: resolve deadbeef: object deadbeef: not found
```

### 3.4 CID 稳定与变化

同一个目录连续导入两次，根 CID 完全相同（内容没变）；只修改 `docs/report.txt`，
则 `report.txt`、`docs`、根节点的 CID 依次变化，而 `docs/notes.txt` 的 CID 不变。
这正是 Merkle DAG 的“变更传播、局部影响”特性。完整演示见
[架构设计说明](docs/架构设计.md#五cid-传播与去重)。

---

## 四、核心原理

### 4.1 对象模型

```go
type Object struct {
    Type  ObjectType `json:"type"`            // blob / tree / list
    Data  []byte     `json:"data,omitempty"`  // Blob 的文件内容
    Links []Link     `json:"links,omitempty"` // Tree/List 的子项
}

type Link struct {
    Name string `json:"name,omitempty"` // 文件名或子目录名
    CID  string `json:"cid"`            // 子对象 CID
    Size int64  `json:"size,omitempty"` // 文件字节数
}
```

### 4.2 CID 计算公式

```text
CID = lowercase_hex( SHA-256( JSON(Object) ) )
```

JSON 由 Go 的 `json.Marshal` 生成：字段顺序固定（Type → Data → Links）、
`omitempty` 规则固定、`[]byte` 按 base64 编码。只要对象相同，字节序列就相同，
CID 就相同——这是“内容寻址”的全部秘密。

### 4.3 Merkle DAG 结构

```text
                    root Tree (ae033b4c...)
                   /          |           \
          README.md        docs Tree      data Tree
        (cfbefcec...)   (1a13f757...)  (b6559387...)
                        /          \
               report.txt        notes.txt
             (9f9b877e...)     (04a66e51...)
```

### 4.4 路径解析

`Resolve(root, "/docs/report.txt")` 把路径切成 `["docs", "report.txt"]`，先加载根
Tree 找到 `docs` 的 CID，再加载 docs Tree 找到 `report.txt` 的 CID，最后读取该
Blob 并校验类型。时间复杂度 O(路径深度 × 目录子项数)，与文件总大小无关。

### 4.5 实测 CID 对照表

| 对象 | 类型 | 大小 | CID |
|------|------|------|-----|
| 根目录 `demo` | tree | — | `ae033b4c1ea3a83944ec106b17e2d080c19848e5f110a930b8e60480083b45b5` |
| `README.md` | blob | 64 B | `cfbefcec90027020a6255f9dd6dd67c33c0e4db20eeed5159818448db27aaed4` |
| `docs` | tree | — | `1a13f7573ff7ee2ea801c398e2c539e8bf61f9d3cd5522c006e0c92401a58872` |
| `docs/report.txt` | blob | 31 B | `9f9b877e2a702f761b4bba5a78bfcf382820a7392babd402a2d685bf84173dae` |
| `docs/notes.txt` | blob | 26 B | `04a66e51e4ed8a86fdaa976cf239f1697c4b686979294c5f8e8c9728d3413952` |
| `data` | tree | — | `b65593877e2428a22863083d97e723bfeab675123f90d6745849e0ecd9a497c7` |
| `data/sample.txt` | blob | 17 B | `b63f4397226b874ce22452d9c8b04d621c1ffefac9c2c58c887cba6d0c2df2b3` |

> 上表是 `testdata/demo` 在 LF 行尾下的实测结果，可直接复现。若在 Windows 上以
> CRLF 检出测试数据，CID 会不同——这也说明 CID 对**字节**敏感，而不是对“文本”敏感。

---

## 五、项目结构

```text
mdag_29/
├── cmd/mdag/main.go        # CLI 入口：add / resolve / cat / ls
├── object/object.go        # 对象模型、CID 计算、尺寸计算
├── store/
│   ├── store.go            # Store 接口
│   ├── memory.go           # 内存存储（测试用）
│   └── file.go             # 文件存储：./data/<cid>.json
├── importer/importer.go    # 递归导入文件/目录，支持 List 分块
├── resolver/resolver.go    # 路径解析、文件读取、目录列举
├── mdag_test.go            # 10 个自动化测试
├── testdata/demo/          # 演示数据（README/docs/data）
├── docs/                   # 项目文档（架构、五看三定、分支管理、协作记录）
├── .github/                # Issue / PR 模板与 CI 工作流
├── CONTRIBUTING.md         # 贡献指南
├── CHANGELOG.md            # 版本记录
├── LICENSE                 # MIT
└── README.md               # 本文件
```

---

## 六、测试

```bash
go test -v -count=1 ./...
```

实测结果：`ok mdag 1.596s`，10 个用例全部通过。

| 测试用例 | 验证内容 |
|----------|---------|
| `TestComputeCID_SameContent` | 相同内容产生相同 CID |
| `TestComputeCID_DifferentContent` | 不同内容产生不同 CID |
| `TestImportTwice_SameCID` | 重复导入幂等，根 CID 不变 |
| `TestImportAndResolve_TwoLevel` | 两级路径解析到正确 Blob |
| `TestResolve_NotFound` | 不存在的路径返回错误 |
| `TestCatOnTree_Error` | 对目录执行 cat 返回错误 |
| `TestModifyFile_ChangesCID` | 修改文件后根/父 CID 变化、无关文件不变 |
| `TestLs` | 列目录返回正确子项数 |
| `TestChunkedRead` | 大文件经 List 分块后可还原 |
| `TestDemoDirectory` | `testdata/demo` 端到端可用 |

CI 使用 GitHub Actions（见 `.github/workflows/ci.yml`），在每次 push / PR 时执行
`gofmt` 检查、`go vet` 与全部测试。

---

## 七、文档索引

| 文档 | 内容 |
|------|------|
| [架构设计说明](docs/架构设计.md) | 对象模型、导入/解析流程、设计权衡、扩展点 |
| [“五看三定”战略分析](docs/五看三定分析.md) | 行业/市场/竞争/自己/机会分析，战略控制点、目标与策略 |
| [分支管理与冲突处理](docs/分支管理与冲突处理.md) | 分支模型、真实合并记录、一次冲突的完整复盘 |
| [协作记录](docs/协作记录.md) | Issue 讨论、Star/Watch、同学互评记录 |
| [作业要求对照表](docs/作业要求对照.md) | 逐条对照期末考核要求与仓库证据 |
| [贡献指南](CONTRIBUTING.md) | 开发环境、分支规范、提交规范、PR 流程 |

---

## 八、分支模型与协作

仓库采用简化版 Git Flow：

- `main`：受保护分支，始终保持可构建、测试通过；
- `feature/*`：功能开发；
- `docs/*`：文档类修改；
- `fix/*`：缺陷修复。

所有改动通过 **Pull Request** 合入 `main`，合并时使用 `--no-ff` 保留分支历史。
提交信息遵循 [Conventional Commits](https://www.conventionalcommits.org/)：
`feat:` / `fix:` / `docs:` / `chore:` / `test:` / `refactor:`。

真实的分支创建、合并与冲突解决过程已记录在
[分支管理与冲突处理](docs/分支管理与冲突处理.md) 中，`git log --graph` 可复核。

---

## 九、路线图

- [x] Blob / Tree / List 对象模型与 SHA-256 CID
- [x] 文件存储与内存存储两种 `Store` 实现
- [x] 目录递归导入（确定性排序）与大文件分块
- [x] 路径解析、文件读取、目录列举
- [x] CLI 四个子命令与错误处理
- [x] 10 个自动化测试 + GitHub Actions CI
- [ ] `mdag fsck`：全量校验对象完整性
- [ ] `mdag gc`：清理不可达对象
- [ ] 把 List 分块接入 CLI（`add --chunk-size`）
- [ ] 只读 HTTP 网关：`GET /ipfs/<cid>/<path>`
- [ ] 对象格式版本号与向后兼容策略

---

## 十、许可证

本项目采用 [MIT License](LICENSE) 开源，欢迎学习、修改与再分发，请保留版权声明。

---

## 致谢

- 课程《开源文化运动》提供的开源实践训练；
- 《工程实践：大数据与人工智能应用》课程中对 Merkle DAG 的理论讲解；
- Git、IPFS、DVC 等开源项目在内容寻址方向上的公开设计与文档。

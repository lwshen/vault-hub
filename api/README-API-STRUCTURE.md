# API 文件结构说明

原来的 `api.yaml` 文件已经被拆分为多个模块化的文件，以提高可维护性和可读性。

## 文件结构

```
api/
├── api-main.yaml              # 主OpenAPI规范文件（引用其他文件）
├── api-original.yaml.bak      # 原始单一文件的备份
├── merge-spec.py              # 合并脚本
├── tool.go                    # Go generate指令
├── cfg.yaml                   # oapi-codegen配置
├── generated.go               # 生成的Go代码
├── paths/                     # API路径定义
│   ├── health.yaml           # 健康检查端点
│   ├── auth.yaml             # 认证相关端点
│   ├── user.yaml             # 用户相关端点
│   ├── vault.yaml            # 保险库相关端点
│   ├── audit.yaml            # 审计日志端点
│   └── api-key.yaml          # API密钥端点
└── schemas/                   # 数据模型定义
    ├── common.yaml           # 通用模型
    ├── auth.yaml             # 认证相关模型
    ├── user.yaml             # 用户相关模型
    ├── vault.yaml            # 保险库相关模型
    ├── audit.yaml            # 审计相关模型
    └── api-key.yaml          # API密钥相关模型
```

## 工作流程

1. **编辑拆分的文件**: 修改 `paths/` 和 `schemas/` 目录中的相应文件
2. **生成合并文件**: 运行 `go generate api/tool.go` 会：
   - 执行 `python3 merge-spec.py` 合并所有拆分的文件为单一的 `api.yaml`
   - 运行 `oapi-codegen` 生成Go代码到 `generated.go`

## 优势

- **模块化**: 每个功能模块有独立的文件
- **可维护性**: 更容易找到和修改特定的API定义
- **可读性**: 文件更小，更容易理解
- **协作友好**: 多人可以同时编辑不同的模块而减少冲突
- **向后兼容**: `go generate` 流程保持不变

## 注意事项

- 不要直接编辑生成的 `api.yaml` 文件，它会被自动覆盖
- 编辑拆分的文件后，需要运行 `go generate` 重新生成代码
- 合并脚本会自动处理文件间的引用关系
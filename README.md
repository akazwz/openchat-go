# Openchat AI Backend API

## 项目简介

Openchat AI Backend API 是一个用 Go (Golang) 编写的后端服务，提供以下核心功能：

- 用户登录和注册
- Token 刷新
- AI 聊天功能
- AI 绘画功能

## 配套前端地址: https://github.com/akazwz/openchat-web

此项目旨在为开发者提供一个简单易用的后端学习示例，使他们能够快速使用Golang开发并集成 AI 功能。
如果你觉得这个项目对你有帮助，请不要忘记给一个 ⭐️ Star！

## 功能特色

- **用户认证系统**: 支持用户注册、登录和基于 token 的身份验证。
- **Token 管理**: 实现安全的 token 刷新机制，确保用户会话的持久性和安全性。
- **AI 聊天**: 集成 AI 聊天功能，支持与 AI 模型进行自然语言对话。
- **AI 绘画**: 利用先进的 AI 技术生成艺术画作。

## 技术栈

- **chi**: 轻量级的路由库，用于处理 HTTP 请求。
- **jwt-go**: 用于生成和验证 JSON Web Tokens，确保安全的用户认证。
- **gorm**: 用于操作数据库，灵活支持多种数据库类型。
- **openai**: OpenAI 官方的 Go SDK，用于实现 AI 聊天和绘画功能。
- **aws s3**: 对象存储，用于存储 AI 绘画生成的图片。
- **blurhash**: 用于生成模糊哈希，实现图片模糊加载效果。

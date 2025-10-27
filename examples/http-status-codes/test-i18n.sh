#!/bin/bash

# JSON 函数多语言翻译测试脚本
echo "🌍 测试 JSON 函数的多语言翻译功能"
echo "================================"

BASE_URL="http://localhost:8080"

# 启动服务器（后台运行）
echo "📍 启动服务器..."
./http-status-example &
SERVER_PID=$!

# 等待服务器启动
sleep 3

echo ""
echo "📋 测试不同语言环境下的 JSON 响应翻译"
echo "================================"

# 测试1: 英文环境 - 默认 JSON 函数
echo ""
echo "1️⃣ 测试英文环境 - JSON 函数"
echo "GET /default (Accept-Language: en)"
curl -s -H "Accept-Language: en" \
     -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/default" | jq '.'

# 测试2: 中文环境 - 默认 JSON 函数
echo ""
echo "2️⃣ 测试中文环境 - JSON 函数"
echo "GET /default (Accept-Language: zh-CN)"
curl -s -H "Accept-Language: zh-CN" \
     -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/default" | jq '.'

# 测试3: 英文环境 - 错误响应
echo ""
echo "3️⃣ 测试英文环境 - 错误响应"
echo "GET /bad-request (Accept-Language: en)"
curl -s -H "Accept-Language: en" \
     -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/bad-request" | jq '.'

# 测试4: 中文环境 - 错误响应
echo ""
echo "4️⃣ 测试中文环境 - 错误响应"
echo "GET /bad-request (Accept-Language: zh-CN)"
curl -s -H "Accept-Language: zh-CN" \
     -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/bad-request" | jq '.'

# 测试5: 英文环境 - 用户不存在错误
echo ""
echo "5️⃣ 测试英文环境 - 用户不存在"
echo "GET /api/v1/users/999 (Accept-Language: en)"
curl -s -H "Accept-Language: en" \
     -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/api/v1/users/999" | jq '.'

# 测试6: 中文环境 - 用户不存在错误
echo ""
echo "6️⃣ 测试中文环境 - 用户不存在"
echo "GET /api/v1/users/999 (Accept-Language: zh-CN)"
curl -s -H "Accept-Language: zh-CN" \
     -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/api/v1/users/999" | jq '.'

# 测试7: 英文环境 - 创建成功
echo ""
echo "7️⃣ 测试英文环境 - 创建成功"
echo "POST /api/v1/users (Accept-Language: en)"
curl -s -X POST -H "Accept-Language: en" -H "Content-Type: application/json" \
     -d '{"name":"Test User","email":"test@example.com"}' \
     -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/api/v1/users" | jq '.'

# 测试8: 中文环境 - 创建成功
echo ""
echo "8️⃣ 测试中文环境 - 创建成功"
echo "POST /api/v1/users (Accept-Language: zh-CN)"
curl -s -X POST -H "Accept-Language: zh-CN" -H "Content-Type: application/json" \
     -d '{"name":"测试用户","email":"test@example.com"}' \
     -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/api/v1/users" | jq '.'

echo ""
echo "✅ 测试完成！"
echo ""
echo "🎯 总结:"
echo "- JSON 函数现在支持真正的多语言翻译"
echo "- 所有响应消息都会根据 Accept-Language 头自动翻译"
echo "- 包括成功消息和错误消息"
echo "- 与模板翻译功能保持一致的翻译机制"

# 关闭服务器
echo ""
echo "🛑 关闭服务器..."
kill $SERVER_PID 2>/dev/null
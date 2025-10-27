#!/bin/bash

# HTTP 状态码测试脚本
echo "🚀 测试 HTTP 状态码自定义功能"
echo "================================"

BASE_URL="http://localhost:8080"

# 启动服务器（后台运行）
echo "📍 启动服务器..."
./http-status-example &
SERVER_PID=$!

# 等待服务器启动
sleep 3

echo ""
echo "📋 测试不同的 HTTP 状态码响应"
echo "================================"

# 测试1: 默认状态码 (200)
echo ""
echo "1️⃣ 测试默认状态码 (200)"
echo "GET /default"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/default" | jq .

# 测试2: 创建成功 (201)
echo ""
echo "2️⃣ 测试资源创建成功 (201)"
echo "POST /created"
curl -s -X POST -H "Content-Type: application/json" \
     -d '{"name":"测试用户","email":"test@example.com"}' \
     -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/created" | jq .

# 测试3: 错误请求 (400)
echo ""
echo "3️⃣ 测试错误请求 (400)"
echo "GET /bad-request"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/bad-request" | jq .

# 测试4: 无法处理的实体 (422)
echo ""
echo "4️⃣ 测试无法处理的实体 (422)"
echo "GET /unprocessable"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/unprocessable" | jq .

# 测试5: 请求已接受 (202)
echo ""
echo "5️⃣ 测试请求已接受 (202)"
echo "PUT /accepted"
curl -s -X PUT -H "X-Request-ID: req-123" \
     -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/accepted" | jq .

# 测试6: 模板参数响应 (201)
echo ""
echo "6️⃣ 测试模板参数响应 (201)"
echo "GET /template"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/template" | jq .

# 测试7: 多语言模板响应 - 英文
echo ""
echo "7️⃣ 测试多语言模板响应 - 英文"
echo "GET /template/i18n (英文)"
curl -s -H "Accept-Language: en" -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/template/i18n" | jq .

# 测试8: 多语言模板响应 - 中文
echo ""
echo "8️⃣ 测试多语言模板响应 - 中文"
echo "GET /template/i18n (中文)"
curl -s -H "Accept-Language: zh-CN" -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/template/i18n" | jq .

# 测试9: 多语言错误模板 - 英文
echo ""
echo "9️⃣ 测试多语言错误模板 - 英文"
echo "GET /template/error (英文)"
curl -s -H "Accept-Language: en" -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/template/error" | jq .

# 测试10: 多语言错误模板 - 中文
echo ""
echo "🔟 测试多语言错误模板 - 中文"
echo "GET /template/error (中文)"
curl -s -H "Accept-Language: zh-CN" -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/template/error" | jq .

# 测试11: RESTful API - 获取用户列表
echo ""
echo "1️⃣1️⃣ 测试 RESTful API - 获取用户列表"
echo "GET /api/v1/users"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/api/v1/users" | jq .

# 测试8: RESTful API - 创建用户
echo ""
echo "1️⃣2️⃣ 测试 RESTful API - 创建用户"
echo "POST /api/v1/users"
curl -s -X POST -H "Content-Type: application/json" \
     -d '{"name":"新用户","email":"newuser@example.com"}' \
     -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/api/v1/users" | jq .

# 测试9: RESTful API - 获取不存在的用户
echo ""
echo "9️⃣ 测试 RESTful API - 获取不存在的用户"
echo "GET /api/v1/users/999"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/api/v1/users/999" | jq .

# 测试10: RESTful API - 更新用户
echo ""
echo "🔟 测试 RESTful API - 更新用户"
echo "PUT /api/v1/users/1"
curl -s -X PUT -H "Content-Type: application/json" \
     -d '{"name":"更新后的用户名","email":"updated@example.com"}' \
     -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/api/v1/users/1" | jq .

# 测试11: RESTful API - 删除用户
echo ""
echo "1️⃣1️⃣ 测试 RESTful API - 删除用户"
echo "DELETE /api/v1/users/2"
curl -s -X DELETE -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/api/v1/users/2"

# 测试12: 场景说明
echo ""
echo "1️⃣2️⃣ 测试不同场景的状态码说明"
echo "GET /scenarios"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/scenarios" | jq .[0]

echo ""
echo "✅ 测试完成！"
echo ""
echo "🎯 总结:"
echo "- 成功演示了自定义 HTTP 状态码功能"
echo "- 展示了不同业务场景下的状态码使用"
echo "- 验证了 RESTful API 的标准状态码响应"
echo "- 演示了模板参数动态消息生成功能"
echo "- 实现了真正的多语言翻译功能"
echo "- 支持通过 Accept-Language 头自动语言检测"

# 关闭服务器
echo ""
echo "🛑 关闭服务器..."
kill $SERVER_PID 2>/dev/null
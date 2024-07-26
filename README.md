# 執行 Docker Compose
確保你在此目錄下執行以下命令来啟動服務：

```sh
docker-compose up --build
```

# API 測試指南

以下是如何測試會員註冊、登入和查詢 API 的步驟。

## 1. 註冊用戶

使用 `curl` 命令來測試用戶註冊 API。

### 請求

```sh
curl -X POST http://localhost:8090/api/register \
     -d '{"username": "testuser", "password": "testpass"}' \
     -H "Content-Type: application/json"
```


## 2. 登錄用戶
使用 curl 命令來測試用戶登入 API。

請求

```sh
curl -X POST http://localhost:8090/api/login \
     -d '{"username": "testuser", "password": "testpass"}' \
     -H "Content-Type: application/json"
```


## 3. 查詢用戶
使用 curl 命令來測試查詢用戶 API。

### 請求

```sh
curl -X GET http://localhost:8090/api/user/testuser
```




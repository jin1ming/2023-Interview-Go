# HTTP协议面试题

[TOC]

## 基础篇

### Q1: HTTP 常见的状态码有哪些？
- **2xx 成功**
  - 200 OK：请求成功。
  - 201 Created：资源创建成功。
  - 204 No Content：请求成功但无返回内容（常用于 DELETE）。
- **3xx 重定向**
  - 301 Moved Permanently：永久重定向（浏览器自动缓存）。
  - 302 Found：临时重定向。
  - 304 Not Modified：资源未修改，使用缓存。
- **4xx 客户端错误**
  - 400 Bad Request：请求参数错误。
  - 401 Unauthorized：未授权（需 Token）。
  - 403 Forbidden：禁止访问（有 Token 但无权限）。
  - 404 Not Found：资源不存在。
- **5xx 服务端错误**
  - 500 Internal Server Error：服务器内部错误。
  - 502 Bad Gateway：网关错误（上游服务报错）。
  - 503 Service Unavailable：服务不可用（超载或维护）。
  - 504 Gateway Timeout：网关超时。

### Q2: GET 和 POST 的区别？
1.  **语义**：GET 用于获取资源（幂等、安全）；POST 用于提交/修改资源（非幂等）。
2.  **参数**：GET 参数放在 URL 中（长度受限，不安全）；POST 参数放在 Body 中。
3.  **缓存**：GET 请求会被浏览器主动缓存；POST 默认不会。
4.  **数据包**：GET 产生一个 TCP 包；POST 可能产生两个 TCP 包（先发 Header，服务器响应 100 Continue 后再发 Body，视浏览器/库而定）。

### Q3: HTTP/1.0, 1.1, 2.0, 3.0 的区别？
- **HTTP/1.0**：短连接，每次请求都要重新建立 TCP 连接。
- **HTTP/1.1**：
  - **长连接 (Keep-Alive)**：复用 TCP 连接。
  - **管道化 (Pipelining)**：支持同时发送多个请求，但响应必须按顺序返回（队头阻塞）。
- **HTTP/2.0**：
  - **多路复用 (Multiplexing)**：基于二进制分帧，真正解决应用层队头阻塞，一个连接并发处理多个请求。
  - **头部压缩 (HPACK)**：减少 Header 传输体积。
  - **服务端推送 (Server Push)**。
- **HTTP/3.0**：
  - **基于 QUIC (UDP)**：解决 TCP 层的队头阻塞问题。
  - **0-RTT 建连**：更快的连接建立速度。
  - **连接迁移**：网络切换（如 WiFi -> 4G）时连接不断。

## 进阶篇

### Q4: HTTPS 的工作原理（握手过程）？
1.  **Client Hello**：客户端发送支持的加密套件、随机数 ClientRandom。
2.  **Server Hello**：服务端选择加密套件，发送证书（含公钥）、随机数 ServerRandom。
3.  **验证证书**：客户端验证证书合法性（CA 签名、有效期、域名）。
4.  **预主密钥**：客户端生成 Pre-Master Secret，用服务端公钥加密后发送。
5.  **生成会话密钥**：双方利用 ClientRandom + ServerRandom + Pre-Master Secret 生成对称密钥 (Session Key)。
6.  **加密通信**：后续数据传输均使用 Session Key 进行对称加密。

### Q5: 什么是跨域 (CORS)？如何解决？
- **定义**：浏览器同源策略限制，协议、域名、端口任一不同即为跨域。
- **解决**：服务端设置响应头：
  - `Access-Control-Allow-Origin`: 允许的域名。
  - `Access-Control-Allow-Methods`: 允许的方法 (GET, POST...)。
  - `Access-Control-Allow-Headers`: 允许的 Header。
  - `Access-Control-Allow-Credentials`: 是否允许携带 Cookie。

### Q6: Cookie 和 Session 的区别？
- **存储位置**：Cookie 存放在客户端（浏览器）；Session 存放在服务端（内存/Redis）。
- **安全性**：Cookie 易被篡改/窃取（XSS/CSRF）；Session 相对安全。
- **容量**：Cookie 通常限制 4KB；Session 无限制。
- **关联**：通常 Session ID 存储在 Cookie 中，用于服务端识别用户身份。

### Q7: 常见的 Web 攻击有哪些？
- **XSS (跨站脚本攻击)**：注入恶意脚本。防御：转义输入输出、CSP。
- **CSRF (跨站请求伪造)**：利用用户登录态发恶意请求。防御：CSRF Token、SameSite Cookie。
- **SQL 注入**：防御：预编译 SQL (Prepare Statement)。

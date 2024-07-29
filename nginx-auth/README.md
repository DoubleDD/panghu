## 在web系统中如何给文件增加认证功能

### 静态文件服务（例如：minio）

### 反向代理服务nginx

后端校验服务需要根据请求的文件路径进行校验，并返回特定的HTTP响应。具体来说，后端校验服务需要：

校验通过时返回200状态码：如果文件校验通过，后端校验服务应该返回HTTP状态码200（OK），并且可以返回一些额外的头部信息，比如X-Accel-Redirect头部，用于指示Nginx直接返回文件内容。
校验不通过时返回401状态码：如果文件校验不通过，后端校验服务应该返回HTTP状态码401（Unauthorized）。
以下是一个示例配置，展示了如何使用X-Accel-Redirect头部来实现这个功能：

```nginx
http {
    upstream backend_check {
        server backend-check:3000;
    }

    server {
        listen 80;
        server_name yourdomain.com;

        location / {
            # 首先将请求转发到后端校验服务
            proxy_pass http://backend_check;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;

            # 根据后端校验服务的响应进行处理
            proxy_intercept_errors on;
            error_page 401 = @unauthorized;
        }

        location @unauthorized {
            return 401;
        }

        location /protected/ {
            internal;
            # 使用X-Accel-Redirect头部来指示Nginx直接返回文件内容
            alias /path/to/protected/files/;
        }
    }
}
```

在这个配置中：

1. location / 块将所有请求转发到后端校验服务。
2. proxy_intercept_errors on; 和 error_page 401 = @unauthorized; 用于捕获401状态码并重定向到 @unauthorized 块。
3. location @unauthorized 块返回401状态码。
4. location /protected/ 块是一个内部位置块，用于处理X-Accel-Redirect头部。这个块使用alias指令来指定实际文件的路径。


后端校验服务在返回200状态码时，应该包含一个X-Accel-Redirect头部，其值为文件的相对路径（相对于/protected/路径）。 例如：

```http
HTTP/1.1 200 OK
X-Accel-Redirect: /protected/path/to/file.pdf
```

这样，Nginx会根据`X-Accel-Redirect`头部的值，直接从指定的路径返回文件内容。

> 总结一下，后端校验服务需要根据文件路径进行校验，并在校验通过时返回200状态码和X-Accel-Redirect头部，在校验不通过时返回401状态码。Nginx会根据这些响应进行相应的处理。



## X-Accel-Redirect

`X-Accel-Redirect` 是 Nginx 提供的一种内部重定向机制，用于将请求重定向到内部位置块（internal location），从而实现对静态文件的安全访问控制。它的工作原理如下：

1. 后端校验：当客户端请求一个文件时，Nginx 首先将请求转发到后端校验服务。后端校验服务根据请求的文件路径进行校验，判断用户是否有权限访问该文件。

2. 返回响应：如果校验通过，后端校验服务返回一个 HTTP 200 状态码，并在响应头部中包含 X-Accel-Redirect 头部。这个头部的值是一个内部 URI，指向实际的文件路径。例如：
```http
HTTP/1.1 200 OK
X-Accel-Redirect: /protected/path/to/file.pdf
```

3. 内部重定向：Nginx 收到后端校验服务的响应后，会检查响应头部中的 X-Accel-Redirect 头部。如果存在这个头部，Nginx 会执行一个内部重定向，将请求转发到指定的内部位置块。

4. 内部位置块：在 Nginx 配置中，需要定义一个内部位置块，用于处理 X-Accel-Redirect 头部指定的 URI。这个位置块通常使用 internal 指令标记为内部位置，确保它只能通过内部重定向访问，而不能直接被客户端请求。例如：
```nginx
location /protected/ {
    internal;
    alias /path/to/protected/files/;
}
```

5. 返回文件内容：当 Nginx 将请求重定向到内部位置块时，它会根据 alias 指令指定的路径，直接返回文件内容给客户端。


通过这种方式，X-Accel-Redirect 实现了对静态文件的安全访问控制。后端校验服务负责权限校验，而 Nginx 负责文件的实际传输。这种分离的设计使得权限校验和文件传输可以独立处理，提高了系统的灵活性和安全性。

> 总结一下，X-Accel-Redirect 的工作原理是：后端校验服务进行权限校验，并在校验通过时返回包含 X-Accel-Redirect 头部的响应；Nginx 根据这个头部执行内部重定向，将请求转发到内部位置块，从而返回文件内容给客户端。


## 测试

1. 启动web服务
```bash
./nginx-auth
```

2. 启动nginx服务
```bash
docker compose up
```

3. 访问图片
    1. 直接访问

    http://localhost/auth/test.jpg

    2. 带认证参数`token=123`访问

    http://localhost/auth/test.jpg?token=123




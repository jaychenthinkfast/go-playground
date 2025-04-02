# 第一阶段：构建Vue应用
FROM docker.m.daocloud.io/library/node:18-alpine as build-stage

# 设置工作目录
WORKDIR /app

# 设置环境变量
ENV NODE_ENV=production

# 替换Alpine镜像源为阿里云
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 复制依赖文件
COPY package*.json ./

# 安装依赖
RUN npm ci --quiet || npm install --quiet

# 复制项目文件
COPY . .

# 构建项目
RUN npm run build

# 第二阶段：Nginx服务阶段
FROM docker.m.daocloud.io/library/nginx:stable-alpine as production-stage

# 设置工作目录
WORKDIR /usr/share/nginx/html

# 替换Alpine镜像源为阿里云
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 从构建阶段复制构建文件
COPY --from=build-stage /app/dist .

# 复制Nginx配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 添加健康检查
RUN apk add --no-cache curl
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/ || exit 1

# 暴露端口
EXPOSE 80

# 启动Nginx服务
CMD ["nginx", "-g", "daemon off;"] 
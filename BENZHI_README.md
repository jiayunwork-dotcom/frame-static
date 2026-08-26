# frame-static：Go 平面刚架静力 Web 服务（直接刚度法 + 前端求解台）

提交节点坐标、杆件截面与荷载，用划行划列直接刚度法解节点位移并回代杆端力。全部支座反力与外荷载主矢、主矩代数和为零；轴压会降低侧向刚度（几何刚度），均温升温在两端固支杆上产生一对反向轴力。矩阵奇异返回错误而非 NaN。

## 构建 / 运行 / 测试

```text
go build ./...
./frame-static -http :8080
curl -s http://127.0.0.1:8080/api/meta
go run . example/portal.json
go test ./...
```

## 评测镜像

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -d -P --name frame-static-b14 <image-name>:latest
curl -s http://127.0.0.1:$(docker port frame-static-b14 8080 | cut -d: -f2)/api/meta
docker rm -f frame-static-b14
```

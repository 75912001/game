#!/bin/bash

# 获取当前脚本文件所在路径的绝对路径
currentPath=$(realpath "$(dirname "$0")")
echo "currentPath:${currentPath}"
#项目绝对路径
projectPath=$(dirname "${currentPath}")
echo "projectPath:${projectPath}"

cd "${projectPath}"/proto || exit

# 检查 protoc 是否安装
if ! command -v protoc &> /dev/null; then
    echo "错误: protoc 未安装"
    echo "请先安装 protoc: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# 检查必要的 Go 插件是否安装
if ! command -v protoc-gen-go &> /dev/null; then
    echo "安装 protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "安装 protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# 创建输出目录（如果不存在）
OutputDir=../pb
mkdir -p ${OutputDir}


# 生成 Go 代码
protoc --go_out=${OutputDir} --go_opt=paths=source_relative \
  ./*.proto

# 检查生成是否成功
if [ $? -eq 0 ]; then
    echo "✅ 代码生成成功！"
    echo "生成的文件位于: ${OutputDir}"
else
    echo "❌ 代码生成失败"
    exit 1
fi

# 将生成的文件cp到 sa/client/src/proto/ 中
ClientPbDir="${projectPath}/client/src/proto"
mkdir -p "${ClientPbDir}"
rm -f ${ClientPbDir}/*.pb.go
cp ${OutputDir}/*.pb.go "${ClientPbDir}/"
if [ $? -eq 0 ]; then
    echo "✅ 文件已复制到: ${ClientPbDir}"
else
    echo "⚠️  文件复制失败"
fi

# 清理临时生成的文件夹
rm -rf ${OutputDir}

# 任意键继续
read -n 1 -s -r -p "按任意键继续..."

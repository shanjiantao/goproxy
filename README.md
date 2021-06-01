# GOPROXY-loongson

本项目 fork 自 [goproxyio/goproxy](https://github.com/goproxyio/goproxy.git)


# 该项目解决的问题
go mod 本身可以直接下载，而 goproxy 相当于一个中间层，有了这个中间层，便能提供更多的可能和功能。原本 goproxy 的功能是用来进行代理和私有仓库配置，这两点原项目本身已经实现。

本项目提供额外的第三个功能，以解决在 loongarch64 下编译 go 项目时，以前固定版本的依赖包由于不支持 loongarch64 架构导致编译失败的问题。

# 实现原理
根据问题产生情况，以 golang.org/x/sys 包为例，该包通过 commit 和时间作为版本。截至 2021 年 5月31日，官方的 golang.org/x/sys 仍然不支持 loongarch64 架构。如果有一个项目用到了该包的某一个版本，通过任意一个 mod 源获取的下来包都会因为不支持 loongarch64 无法顺利完成编译。而且考虑到不同的项目所依赖的 golang/x/sys 版本是不一样的，且每个版本要适配 loongarch64 架构的代码也都相同，所以考虑这样一种解决办法。即在 proxy 端一个项目只保留一份固定代码，当 proxy 检测到该项目的任何版本时，都从该固定版本打包返回。该实现方案能解决重复适配的问题，但仍然有以下约束。

# 方案的缺陷
1. 不能使用 go sum 机制，很明显我们的适配破坏了官方的 sum 值。所以请删除项目的 go.sum 文件，并设置该环境变量 `GOSUMDB='off'` 
2. 项目应该满足高版本兼容低版本，因为我们只保留了一个版本的代码。


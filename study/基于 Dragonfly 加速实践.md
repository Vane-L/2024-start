source link: [火山引擎基于 Dragonfly 加速实践](https://www.sofastack.tech/blog/volcano-engine-based-on-dragonfly-acceleration-practices/)

### 背景
- 客户端数量越来越多，镜像越来越大，拉取并发量受限于带宽和QPS。
- 因此需要引入P2P，减轻服务端压力，进而满足大规模并发拉取镜像的需求。



### Dragonfly 
Manager：维护各个P2P集群之间的关系，动态配置管理和RBAC。还包括前端控制台，方便用户对集群进行可视化操作。

Scheduler：为下载节点选择最优的下载父节点。异常控制Dfdaemon的回源。

Seed Peer：Dfdaemon开启Seed Peer模式可以作为P2P集群中的回源下载peer，是整个集群中下载的根peer。

Peer : 与dfdaemon一起部署，基于C/S架构，提供dfget命令下载工具，以及dfget daemon运行守护进程提供任务下载能力。

### Kraken
- 代理人
    - 部署在每个主机上
    - 实现Docker注册接口
    - 向追踪器公布可用内容
    - 连接到跟踪器返回的对等点以下载内容
- 起源
    - 专用播种机
    - 将 blob 作为文件存储在由可插拔存储（例如 S3、GCS、ECR）支持的磁盘上
    - 形成一个自我修复的哈希环来分配负载
- 追踪器
    - 跟踪哪些同行拥有哪些内容（正在进行和已完成）
    - 为任意给定的 blob 提供要连接的对等点的有序列表
- 代理人
    - 实现Docker注册接口
    - 将每个图像层上传到负责的来源（记住，来源形成一个哈希环）
    - 上传标签以构建索引
- 建立索引
    - 将人类可读的标签映射到 blob 摘要
    - 没有一致性保证：客户端应该使用唯一的标签
    - 支持集群之间的图像复制（带重试的简单重复队列）
    - 将标签作为文件存储在由可插拔存储（例如 S3、GCS、ECR）支持的磁盘上
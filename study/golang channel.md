## 设计原理
Go 语言中最常见的、也是经常被人提及的设计模式就是：
- 不要通过共享内存的方式进行通信，而是应该通过通信的方式共享内存。

目前的 Channel 收发操作均遵循了先进先出的设计，具体规则如下：
1. 先从 Channel 读取数据的 Goroutine 会先接收到数据；
2. 先向 Channel 发送数据的 Goroutine 会得到先发送数据的权利；

## 发送数据
使用 ch <- i 表达式向 Channel 发送数据时遇到的几种情况：
1. 如果当前 Channel 的 recvq 上存在已经被阻塞的 Goroutine，那么会直接将数据发送给当前 Goroutine 并将其设置成下一个运行的 Goroutine；
2. 如果 Channel 存在缓冲区并且其中还有空闲的容量，我们会直接将数据存储到缓冲区 sendx 所在的位置上；
3. 如果不满足上面的两种情况，会创建一个 runtime.sudog 结构并将其加入 Channel 的 sendq 队列中，当前 Goroutine 也会陷入阻塞等待其他的协程从 Channel 接收数据；

## 接收数据
从 Channel 中接收数据时可能会发生的五种情况：
1. 如果 Channel 为空，那么会直接调用 runtime.gopark 挂起当前 Goroutine；
2. 如果 Channel 已经关闭并且缓冲区没有任何数据，runtime.chanrecv 会直接返回；
3. 如果 Channel 的 sendq 队列中存在挂起的 Goroutine，会将 recvx 索引所在的数据拷贝到接收变量所在的内存空间上并将 sendq 队列中 Goroutine 的数据拷贝到缓冲区；
4. 如果 Channel 的缓冲区中包含数据，那么直接读取 recvx 索引对应的数据；
5. 在默认情况下会挂起当前的 Goroutine，将 runtime.sudog 结构加入 recvq 队列并陷入休眠等待调度器的唤醒；
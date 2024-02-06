source link: [聊聊 controller-runtime 缓存那些事](https://cloud.tencent.com/developer/article/2256963)

### DelegatingClient 对象
- CacheReader
- ClientReader
```go

func NewDelegatingClient(in NewDelegatingClientInput) (Client, error) {
	uncachedGVKs := map[schema.GroupVersionKind]struct{}{}
	for _, obj := range in.UncachedObjects {
		gvk, err := apiutil.GVKForObject(obj, in.Client.Scheme())
		if err != nil {
			return nil, err
		}
		uncachedGVKs[gvk] = struct{}{}
	}

	return &delegatingClient{
		scheme: in.Client.Scheme(),
		mapper: in.Client.RESTMapper(),
		Reader: &delegatingReader{
			CacheReader:       in.CacheReader,
			ClientReader:      in.Client,
			scheme:            in.Client.Scheme(),
			uncachedGVKs:      uncachedGVKs,
			cacheUnstructured: in.CacheUnstructured,
		},
		Writer:                       in.Client,
		StatusClient:                 in.Client,
		SubResourceClientConstructor: in.Client,
	}, nil
}
```
### Get 方法实现
- 是否有 cache object
- 否：直连 api-server get
- 是：直接从 cache get
```go
func (d *delegatingReader) Get(ctx context.Context, key ObjectKey, obj Object, opts ...GetOption) error {
	if isUncached, err := d.shouldBypassCache(obj); err != nil {
		return err
	} else if isUncached {
		return d.ClientReader.Get(ctx, key, obj, opts...)
	}
	return d.CacheReader.Get(ctx, key, obj, opts...)
}
```
### Cache Get
```go
// Get implements Reader.
func (ip *informerCache) Get(ctx context.Context, key client.ObjectKey, out client.Object, opts ...client.GetOption) error {
	gvk, err := apiutil.GVKForObject(out, ip.Scheme)
	if err != nil {
		return err
	}

	started, cache, err := ip.InformersMap.Get(ctx, gvk, out)
	if err != nil {
		return err
	}

	if !started {
		return &ErrCacheNotStarted{}
	}
	return cache.Reader.Get(ctx, key, out)
}
```
### InformersMap.Get
```go
func (ip *specificInformersMap) Get(ctx context.Context, gvk schema.GroupVersionKind, obj runtime.Object) (bool, *MapEntry, error) {
	// Return the informer if it is found
	i, started, ok := func() (*MapEntry, bool, bool) {
		ip.mu.RLock()
		defer ip.mu.RUnlock()
		i, ok := ip.informersByGVK[gvk]
		return i, ip.started, ok
	}()
    // if not found, create new informer
	if !ok {
		var err error
		if i, started, err = ip.addInformerToMap(gvk, obj); err != nil {
			return started, nil, err
		}
	}
    // wait for informer synced
	if started && !i.Informer.HasSynced() {
		// Wait for it to sync before returning the Informer so that folks don't read from a stale cache.
		if !cache.WaitForCacheSync(ctx.Done(), i.Informer.HasSynced) {
			return started, nil, apierrors.NewTimeoutError(fmt.Sprintf("failed waiting for %T Informer to sync", obj), 0)
		}
	}

	return started, i, nil
}
```
### WaitForCacheSync
```go
func WaitForCacheSync(stopCh <-chan struct{}, cacheSyncs ...InformerSynced) bool {
	err := wait.PollImmediateUntil(syncedPollPeriod,
		func() (bool, error) {
			for _, syncFunc := range cacheSyncs {
				if !syncFunc() {
					return false, nil
				}
			}
			return true, nil
		},
		stopCh)
	if err != nil {
		klog.V(2).Infof("stop requested")
		return false
	}

	klog.V(4).Infof("caches populated")
	return true
}
```
### Summary
1. 为开启缓存的资源准备好 ServiceAccount 权限
2. 通常使用 Cache Client 是为了降低控制器 Get/List 请求对 api-server 的负载压力，但这必然会占用更多内存空间

### Q & A
Q: ListWatch 中被 selector 过滤掉的对象一定不会触发对这些对象的 Reconcile 吗？
- 答案：不一定。某类资源的 WorkQueue 消息来源不仅仅是该资源的 Informer，也可以来自其他类型资源的 Informer。
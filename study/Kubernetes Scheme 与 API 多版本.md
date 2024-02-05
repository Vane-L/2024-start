source link: [Kubernetes Scheme 与 API 多版本](https://kayn.wang/kubernetes-scheme/#%E5%AD%98%E5%82%A8%E7%89%88%E6%9C%AC%E5%92%8C%E2%80%9C%E5%86%85%E5%AD%98%E2%80%9D%E7%89%88%E6%9C%AC)
## 概念
### GVK 
> Group and Version and Kinds
- 例子 /apis/batch/v1/namespaces/$NAMESPACE/jobs
  - batch 是 Group
  - v1 是 Version
  - jobs 是 Resource

### Scheme
为了同时解决数据对象的序列化、反序列化与多版本数据对象的兼容和转换问题, Kubernetes设计了一套复杂的机制.
```go
type Scheme struct {
  // gvkToType allows one to figure out the go type of an object with
  // the given version and name.
  gvkToType map[schema.GroupVersionKind]reflect.Type
  // typeToGVK allows one to find metadata for a given go object.
  // The reflect.Type we index by should *not* be a pointer.
  typeToGVK map[reflect.Type][]schema.GroupVersionKind
  // unversionedTypes are transformed without conversion in ConvertToVersion.
  unversionedTypes map[reflect.Type]schema.GroupVersionKind
  // unversionedKinds are the names of kinds that can be created in the context of any group
  // or version
  unversionedKinds map[string]reflect.Type
  // Map from version and resource to the corresponding func to convert
  // resource field labels in that version to internal version.
  fieldLabelConversionFuncs map[schema.GroupVersionKind]FieldLabelConversionFunc
  // defaulterFuncs is a map to funcs to be called with an object to provide defaulting
  // the provided object must be a pointer.
  defaulterFuncs map[reflect.Type]func(interface{})
  // converter stores all registered conversion functions. It also has
  // default converting behavior.
  converter *conversion.Converter
  // versionPriority is a map of groups to ordered lists of versions for those groups indicating the
  // default priorities of these versions as registered in the scheme
  versionPriority map[string][]string
  // observedVersions keeps track of the order we've seen versions during type registration
  observedVersions []schema.GroupVersion
  schemeName string
}
```
Scheme 资源类型的注册
```go
// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: "apps", Version: "v1beta1"}

var (
  SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
  localSchemeBuilder = &SchemeBuilder
  AddToScheme        = localSchemeBuilder.AddToScheme
)

// Adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
  scheme.AddKnownTypes(SchemeGroupVersion,
    &Deployment{},
    &DeploymentList{},
    &DeploymentRollback{},
    &Scale{},
    &StatefulSet{},
    &StatefulSetList{},
    &ControllerRevision{},
    &ControllerRevisionList{},
  )
  metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
  return nil
}
```
Scheme Defaulter
```go
func SetDefaults_StatefulSet(obj *appsv1beta1.StatefulSet) {
  ...
  labels := obj.Spec.Template.Labels
  if labels != nil {
    if obj.Spec.Selector == nil {
      obj.Spec.Selector = &metav1.LabelSelector{
        MatchLabels: labels,
      }
    }
    if len(obj.Labels) == 0 {
      obj.Labels = labels
    }
  }
  ...
}
```
Scheme Converter
```go
// AddConversionFunc registers a function that converts between a and b by passing objects of those
// types to the provided function. The function *must* accept objects of a and b - this machinery will not enforce
// any other guarantee.
func (s *Scheme) AddConversionFunc(a, b interface{}, fn conversion.ConversionFunc) error {
  return s.converter.RegisterUntypedConversionFunc(a, b, fn)
}

// AddGeneratedConversionFunc registers a function that converts between a and b by passing objects of those
// types to the provided function. The function *must* accept objects of a and b - this machinery will not enforce
// any other guarantee.
func (s *Scheme) AddGeneratedConversionFunc(a, b interface{}, fn conversion.ConversionFunc) error {
  return s.converter.RegisterGeneratedUntypedConversionFunc(a, b, fn)
}
```
## 使用 Scheme
### 生成静态 Client
通过 https://github.com/kubernetes/code-generator 生成 Client 代码, 项目引用 Client 后进行相关操作
```go
var Scheme = runtime.NewScheme()
var Codecs = serializer.NewCodecFactory(Scheme)
var ParameterCodec = runtime.NewParameterCodec(Scheme)
var localSchemeBuilder = runtime.SchemeBuilder{
// 此处注册资源的Scheme
  xxx.AddToScheme,
}
var AddToScheme = localSchemeBuilder.AddToScheme

func init() {
  v1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
  utilruntime.Must(AddToScheme(Scheme))
}
```
### 使用动态 Client
唯一需要的一般是在程序初始化时将相关资源的 Scheme 注册到 controller-runtime Client 的 Scheme 中
```go
err := client.Get(ctx, req.NamespacedName, &obj)
```

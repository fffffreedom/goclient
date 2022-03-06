# client-go 实践

## 环境
```
$ go version
go version go1.16 darwin/amd64
```
```
$ kubectl version
Client Version: version.Info{Major:"1", Minor:"17", GitVersion:"v1.17.3", GitCommit:"06ad960bfd03b39c8310aaf92d1e7c12ce618213", GitTreeState:"clean", BuildDate:"2020-02-11T18:14:22Z", GoVersion:"go1.13.6", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"17", GitVersion:"v1.17.3", GitCommit:"06ad960bfd03b39c8310aaf92d1e7c12ce618213", GitTreeState:"clean", BuildDate:"2020-02-11T18:07:13Z", GoVersion:"go1.13.6", Compiler:"gc", Platform:"linux/amd64"}
```

## 问题

### 问题一
执行go mod init/tidy/vendor后，运行go run out-cluster-client-configuration.go，报错如下：
```
../vendor/sigs.k8s.io/json/internal/golang/encoding/json/encode.go:1249:12: sf.IsExported undefined (type reflect.StructField has no field or method IsExported)
../vendor/sigs.k8s.io/json/internal/golang/encoding/json/encode.go:1255:18: sf.IsExported undefined (type reflect.StructField has no field or method IsExported)
```
参考如下，升级一下go版本即可：
https://github.com/clarketm/json/issues/5

### 问题二
in-cluster-client-configuration测试问题
1、使用scratch基础镜像
由于镜像的启动命令配置为ENYRYPOINT app，启动时命令会转换成：ENTRYPOINT ["/bin/sh","-c","echo ..."]，
也就是会用到sh，但scratch为空镜像，没有shell环境，故在启动时会失败。
https://www.codenong.com/54820846/

2、权限问题
镜像启动后，会报权限错误，即对于default namespace，没有权限查看
If you have RBAC enabled on your cluster, use the following snippet to create role binding 
which will grant the default service account view permissions.
kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default

3、问题排查
kubectl logs pod-name
kubectl run 使用-i参数，让pod的输出直接输出到终端
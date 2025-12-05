一、运维系列教程
运维手册
二、通用
1. rtsp录制与流媒体服务器推流
1. 启动流媒体服务 docker run --rm -it -e MTX_PROTOCOLS=tcp -p 8554:8554 -p 1935:1935 -p 8888:8888 -p 8889:8889 aler9/rtsp-simple-server
2. 录制相机视频流 ffmpeg -i rtsp://admin:tman1234@192.168.30.64/h264/main/av_stream -vcodec copy -acodec copy output.mp4 （如果已经存在视频文件，可跳过此步骤）
3. 循环推流 ffmpeg -re -stream_loop -1 -i output.mp4 -vcodec copy -acodec copy -f flv rtmp://127.0.0.1:1935/1 
4. 打开浏览器http://192.168.30.14:8888/1/
2. 查看服务器出口 IP 命令
curl cip.cc
3. 挂载腾讯云对象存储到服务器
参考文档：
1. https://cloud.tencent.com/document/product/436/6883
2. https://cloud.tencent.com/developer/article/1855290?from=15425
4. 清理 CI 机器镜像
docker system prune --volumes

docker rmi $(docker images| awk '{print $1":"$2}')

    CONTAINERD_ADDRESS="unix:///run/k3s/containerd/containerd.sock"
    
    nerdctl image prune
5. CD HTTPS 证书更新异常
大概率原因是因为自动更新证书的组件出问题了，可以检查 cert-manager-webhook 应用是否正常，检查 cert-manager-webhook-ca 证书是否过期。如果过期，直接删除证书即可，后续会自动新建。
# 查看证书过期时间
kubectl get secret cert-manager-webhook-ca -n cert-manager -o jsonpath="{.data['ca\.crt']}" | base64 --decode > cert-manager-webhook-ca.crt

openssl x509 -in cert-manager-webhook-ca.crt -text -noout | grep 'Not After'

# print :  Not After : Oct 24 23:01:57 2024 GMT

# 备份证书
kubectl get secret cert-manager-webhook-ca -n cert-manager -o yaml > cert-manager-webhook-ca-backup.yaml

# 删除证书
kubectl delete secret -n cert-manager cert-manager-webhook-ca

# 重启应用
kubectl rollout restart deployment cert-manager-webhook -n cert-manager
6. Git lfs 同步异常的处理方法
本质上是，从源仓库把所有 lfs 文件拉下来，然后手动 push 到 mirror 仓库里
# 进入CI 服务器
cd /home/ubuntu/wangtao/motion/motion
# 确认git remote 地址是  git@gitlab.deepglint.com:alpha/motion.git
cat .git/config | more
git fetch
git branch -r | grep -v '\->' | while read remote; do git branch --track "${remote#origin/}" "$remote"; done
git lfs fetch --all
git lfs pull
# 更换git remote源
# http://root:deep2013@10.0.0.24:32745/alpha/motion.git
vim .git/config
git lfs push origin --all
三、DIPPER
1. 如何对边缘服务器批量执行 bash 命令？
【Dipper】边缘服务器批量执行 Bash 命令脚本
2. Ubuntu 扩展逻辑分区
https://cloud.tencent.com/developer/article/1965711
lvdisplay
lsblk
lvextend -l +100%FREE /dev/ubuntu-vg/ubuntu-lv
resize2fs /dev/ubuntu-vg/ubuntu-lv
lvdisplay
lsblk
3. dipper 升级
helmrelease 都应该为deployed或则为空，显示failed都是有问题，可以看一下具体问题，修复后删除helmrelease重新升级一下
4. 临时修改 vse 授权地址用于测试
 
kubectl -n vse get cm vse-config -o yaml |  sed -e 's|pastur.pastur.svc.cluster.local.|192.168.3.79:2243|' |  sed -e 's|"LicenseType":"local"|"LicenseType":"concurrency"|' | kubectl apply -f -

kubectl -n vse delete po -l app=vse
 
# 还原 vse 授权配置
kubectl -n vse get cm vse-config -o yaml |  sed -e 's|192.168.3.79:2243|pastur.vse.svc.cluster.local.|' |  sed -e 's|"LicenseType":"concurrency"|"LicenseType":"local"|' |  kubectl apply -f -

kubectl -n vse delete po -l app=vse
5. 远程到边缘服务器的方法
以下是三种方式，不是三步操作
1. 普通方法：gitlab.deepglint.com
如果需要提供其他同事远程运维通道，或者 vse Pod 应用异常时，可以采用下面这种方式
2. 结合上述方法1，代理 nginx Pod 22 端口到本地主机某一端口（如下图）；以此类推可以代理边缘机器的任意端口到本地。
[图片]
3. 终极方法 远程访问机器教程 （使用 nps ）
6. 查看cuda核心数
在宿主机可直接执行
wget https://star-deploy-1253924368.cos.ap-beijing.myqcloud.com/tools/cuda_info
chmod +x cuda_info
./cuda_info
7. 删除namespace下异常pod
(
for ns in $(kubectl get ns)
do
kubectl -n $ns delete po $(kubectl get po -n $ns | grep "ContainerStatusUnknown\|Terminating\|Evicted\|Error\|Completed\|CrashLoopBackOff" | awk '{print $1}')
done
)
删除失败的helmrelease
(
for ns in $(kubectl get ns)
do
kubectl -n $ns delete helmrelease --force --grace-period=0  $(kubectl get helmrelease -n $ns | grep "pending-upgrade" | awk '{print $1}')
done
)
8. helmRelease丢失
以机器id 210235a2cr522cc0047t来举例
# 先去机器上，kubectl scale deploy clusternet-agent -n clusternet-system --replicas 0

kubectl get mcls -A | grep 210235a2cr522cc0047t
210235a2cr522cc0047t   210235a2cr522cc0047t   8c44987b-332d-474d-8125-ca9c0cf908ee   Dual        v1.26.4+k3s1    True      2m6s
# 第三列是它的cluster_id

kubectl get clsrr -A | grep 8c44987b-332d-474d-8125-ca9c0cf908ee
clusternet-8c44987b-332d-474d-8125-ca9c0cf908ee   8c44987b-332d-474d-8125-ca9c0cf908ee   Approved   3s
# 第一列是它的凭证名称

kubectl delete clsrr clusternet-8c44987b-332d-474d-8125-ca9c0cf908ee
kubectl delete namespace 210235a2cr522cc0047t
# 删除clsrr凭证和机器的namespace

# 再回到机器上，kubectl scale deploy cluster net-agent -n clusternet-system --replicas 1

9. k3s证书更新
https://docs.k3s.io/cli/certificate#rotating-client-and-server-certificates

检查当前证书过期时间
cat /etc/rancher/k3s/k3s.yaml |grep 'client-certificate-data' | awk '{print $2}' | base64 -d | openssl x509 -text -noout
刷新client和sever证书
systemctl stop k3s

k3s certificate rotate --data-dir /data/k3s

## 检查sever证书过期时间
for i in `ls /data/k3s/server/tls/*.crt`; do echo $i; openssl x509 -enddate -noout -in $i;done

## 重启后 /etc/rancher/k3s/k3s.yaml 证书会自动更新
systemctl start k3s

- k3s 默认安装路径
for i in `ls /var/lib/rancher/k3s/server/tls/*.crt`; do echo $i; openssl x509 -enddate -noout -in $i;done


10. 授权异常
日志中显示有正常的授权范围，但接口返回的是没授权
systemctl restart hasplmd
11. Clusternet 父集群 kubectl 使用异常
问题：
[图片]
解决办法：
# 登陆父集群
rm /var/lib/rancher/k3s/server/tls/dynamic-cert.json
kubectl --server=https://localhost:6444 delete secret -n kube-system k3s-serving

# 然后在数据库 k3s 中，删除所有条件为 name='/registry/secrets/kube-system/k3s-serving' 的 kine 表记录 

systemctl restart k3s
四、STAR
1. Kafka脏数据删除
https://blog.csdn.net/weixin_39992480/article/details/110296521
# Detach Kafka引擎表 ex. audit_queue
detach table name 

# 查看当前的offset, 在我们的开发环境 server 为： 10.0.0.112:30001
 ./kafka-consumer-groups.sh --bootstrap-server 10.0.0.112:9092 --group clickhouse-realtime  --describe 
  
# 根据offset选择指定的new offset
 ./kafka-consumer-groups.sh --bootstrap-server 10.0.0.112:9092 --group clickhouse-realtime  --topic realtime-capture:0 --reset-offsets --to-offset
  807262342 --execute  
  
  # attach Kafka引擎表
  attach table name  
2. kafka跳过脏数据
当重置kafka offset会导致正常的消息丢失时（有多条脏数据，中间夹杂着正常数据），无法通过上述方式跳过，需要设置 kafka engine表, kafka_handle_error_mode= 'stream' ，让clickhouse跳过异常数据。等脏数据处理完成后，再恢复。
- Drop table demo_queue
- 重新创建kafka engine表，添加 kafka_handle_error_mode= 'stream'配置
- 脏数据处理完成后，再次恢复原有的建表语句（把kafka_handle_error_mode= 'stream'删掉）
  - Drop table demo_queue
  - 重新创建kafka engine表，不要添加 kafka_handle_error_mode= 'stream'配置
3. 强制删除namespace(修改代码里的NAMESPACE)
(
NAMESPACE=szdzjs20250200001
kubectl proxy &
kubectl get namespace $NAMESPACE -o json |jq '.spec = {"finalizers":[]}' >temp.json
curl -k -H "Content-Type: application/json" -X PUT --data-binary @temp.json 127.0.0.1:8001/api/v1/namespaces/$NAMESPACE/finalize
)
4. 腾讯云磁盘扩容
在扩容之前把应用先停掉，然后执行如下命令，storage要求必须是10的整数倍
kubectl patch pvc data-postgresql-0 -p '{"spec":{"resources":{"requests":{"storage":"70Gi"}}}}' -n postgresql
5. 重跑冒烟测试
kubectl -n star get job "init-api-smoke-test" -o json | jq 'del(.spec.selector)' | jq 'del(.spec.template.metadata.labels)' | kubectl replace --force -f -
6. 更新集群 Agent Node k3s 证书步骤
1. 首先查看 Node 证书到期时间
for cert in /var/lib/rancher/k3s/agent/*.crt; do
  echo "Checking $cert"
  openssl x509 -in "$cert" -noout -text | grep "Not After"
done
2. 如果发现有证书过期，先停止 k3s 服务：sudo systemctl stop k3s-agent
3. 然后删除过期证书对应的文件，例如（删除前先备份一下）：
sudo rm -rf /var/lib/rancher/k3s/agent/client-k3s-controller.crt
sudo rm -rf /var/lib/rancher/k3s/agent/client-k3s-controller.key
sudo rm -rf /var/lib/rancher/k3s/agent/client-kube-proxy.crt
sudo rm -rf /var/lib/rancher/k3s/agent/client-kube-proxy.key
4. 然后启动 k3s 服务：sudo systemctl start k3s-agent
5. 检查一下证书的最新到期时间
7. 所有边缘设备无法上报属性（可以获取云端属性，但无法上报数据）
其中一个原因是 thingsboard 在监听到 iot.smart.nemorace.com 证书更新后，会自动重启，此时可能会导致 Queue 出问题，全都是 timeout 超时消息，暂时重启 thingsboard 可以解决。
[图片]
五、体育
1. 如何使用 adb 命令调整竖屏音量？
# 实现方式为模拟安卓按键

adb shell input keyevent 24 # 增加音量

adb shell input keyevent 25 # 减小音量

adb shell input keyevent 164 # 静音


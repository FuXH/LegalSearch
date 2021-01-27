# 一、检索服务
## 1、安装jdk
要求jdk1.8版本以上
```shell
java -version
```

## 2、安装elasticsearch
a) yum源
```shell
vim /etc/yum.repos.d/elasticsearch.repo
```
```yaml
[elasticsearch]
name=Elasticsearch repository for 7.x packages
baseurl=https://artifacts.elastic.co/packages/7.x/yum
gpgcheck=1
gpgkey=https://artifacts.elastic.co/GPG-KEY-elasticsearch
enabled=0
autorefresh=1
type=rpm-md
```
```shell
sudo yum install --enablerepo=elasticsearch elasticsearch
```
b) 配置elasticsearch
```shell
vim /etc/elasticsearch/elasticsearch.yml
```
```yaml
network.host: 127.0.0.1
http.port: 9200
```
c) 创建新用户
elasticsearch默认不允许使用root账号启动，user:es, passwd:es
```shell
adduser es
passwd es

chown -R es /usr/local/elasticsearch
```
d) 启动elasticsearch
```shell
systemctl start elasticsearch
```

## 3、部署服务
a) 编译go服务，上传服务器，/usr/local/elasticsearch/service
b) 设置服务配置文件
```shell
vim /usr/local/elasticsearch/service/conf.yaml
```
```yaml
server:
  port: 12341

elasticsearch:
  address:
    - http://127.0.0.1:9200
```
c) 启动服务
```shell
nohup /usr/local/elasticsearch/service/LegalSearch >> service.log &
```
d) 导入数据
```shell
# 清除原有数据
curl 127.0.0.1:8082/api/operation\?action=clean
# 重启服务
nohup /usr/local/elasticsearch/service/LegalSearch >> service.log &
# 导入新数据，path为具体路径
curl 127.0.0.1:8082/api/update\?path=/mnt/data/wenshu
```

# 二、kibana
## 1、安装kibana
yum源
```shell
vim /etc/yum.repos.d/kibana.repo
```
```yaml
[kibana-7.x]
name=Kibana repository for 7.x packages
baseurl=https://artifacts.elastic.co/packages/7.x/yum
gpgcheck=1
gpgkey=https://artifacts.elastic.co/GPG-KEY-elasticsearch
enabled=1
autorefresh=1
type=rpm-md
```
安装kibana
```shell
sudo yum install kibana
```
## 2、配置&启动kibana
a) 配置kibana
```shell
vim /etc/kibana/kibana.yaml
```
```yaml
server.port 5602
server.host "0.0.0.0"
```
b) 启动
```shell
systemctl start kibana
```

## 3、密码登录
a) 生成密码文件
```shell
yum install httpd-tools
mkdir passwd
cd passwd
htpasswd -c -b /etc/nginx/passwd/kibana.passwd test test
```
b) 安装nginx
```shell
yum install nginx

vim /etc/nginx/conf.d/default.conf
```
文件内容
```shell
server {
  listen 5601 default_server;
  
  location / {
    proxy_pass http://127.0.0.1:5602$request_uri;
    
    auth_basic "登录验证";
    auth_basic_user_file /etc/nginx/passwd/kibana.passwd
  }
}
```
c) 启动nginx
```shell
systemcrl start nginx.service
```

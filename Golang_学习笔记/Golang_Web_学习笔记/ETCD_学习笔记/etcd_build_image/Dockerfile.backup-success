FROM centos

MAINTAINER chris

ENV TZ="Asia/Shanghai";\
    DOWNLOAD=https://github.com/etcd-io/etcd/releases/download/v3.4.16/ \
	ETCDVERSION=etcd-v3.4.16-linux-amd64 \
	USER=admin

RUN yum install curl wget tar -y ;\
    useradd ${USER} ;\
    mkdir -p /export/{servers,Logs,packages,Apps,Shell} ;\
    wget ${DOWNLOAD}${ETCDVERSION}.tar.gz && tar -zxf ${ETCDVERSION}.tar.gz -C /export/servers/ && \
    /bin/rm -rf ${ETCDVERSION}.tar.gz ;\
    chown -R ${USER}.${USER} /export ;\
    ln -s /export/servers/${ETCDVERSION}/etcd* /usr/bin/;

EXPOSE 2379 2380

    


FROM ubuntu:20.04

RUN apt-get update -y && apt-get install make openssh-client curl telnet supervisor libaio1 -y

RUN DEBIAN_FRONTEND=noninteractive TZ=Asia/Ho_Chi_Minh apt -y install tzdata

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
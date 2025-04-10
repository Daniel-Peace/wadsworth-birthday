# setting base container image
FROM debian:12.10

# setting user to root
USER root

WORKDIR /app

# updating package lists
RUN apt update

# installing ssh server
RUN apt-get install -y openssh-server
RUN mkdir /var/run/sshd
RUN echo 'root:1c6H0kX#&m' | chpasswd
RUN echo "PermitRootLogin yes" >>  /etc/ssh/sshd_config
COPY ./id_rsa.pub /root/.ssh/authorized_keys

# installing mongodb
RUN apt-get install -y gnupg curl
RUN curl -fsSL https://www.mongodb.org/static/pgp/server-8.0.asc | gpg -o /usr/share/keyrings/mongodb-server-8.0.gpg --dearmor
RUN echo "deb [ signed-by=/usr/share/keyrings/mongodb-server-8.0.gpg ] http://repo.mongodb.org/apt/debian bookworm/mongodb-org/8.0 main" | tee /etc/apt/sources.list.d/mongodb-org-8.0.list
RUN apt-get update
RUN apt-get install -y mongodb-org

# exposing ssh port
EXPOSE 22
EXPOSE 27017

# copying over script to start ssh server and mongodb
COPY ./establish_entrypoints.sh /establish_entrypoints.sh
RUN chmod +x /establish_entrypoints.sh

# starting ssh server and mongodb
ENTRYPOINT ["/establish_entrypoints.sh"]

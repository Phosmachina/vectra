# Use an official Ubuntu as a base image
FROM ubuntu:latest

RUN apt-get update

# Install essential packages
RUN apt-get install --no-install-recommends -y \
    curl wget git openssh-server rsync \
    nodejs npm \
    minify

# Install golang
RUN curl -O https://dl.google.com/go/go1.21.3.linux-amd64.tar.gz  && \
    tar -xvf go1.21.3.linux-amd64.tar.gz && \
    mv go /usr/local

# Set Go and Node.js environment variables
ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
ENV NVM_DIR="/root/.nvm"

# Install autoprefixer
RUN npm install -g postcss postcss-cli autoprefixer

# Install sass
RUN npm install -g sass

# Install jade
RUN go install github.com/Joker/jade/cmd/jade@latest

# Prepare SSH server
RUN mkdir /var/run/sshd
RUN echo 'root:root' | chpasswd
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
EXPOSE 22

# Start SSH server
CMD ["/usr/sbin/sshd", "-D"]

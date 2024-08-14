FROM openjdk:24-jdk-bookworm

ENV SW_AGENT_NAME="oom-test"
ENV TZ="Asia/Shanghai"
ENV HOME="/apps"
ENV JAVA_AGENT_OPTS=$JAVA_AGENT_OPTS
ENV JAVA_OPTS="-Xlog:gc:/apps/jvm-logs/gc.log"
ENV OOM_OPS="-XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/apps/jvm-logs/heapdump.hprof -XX:OnOutOfMemoryError=\"/apps/upload.sh\""


RUN mkdir -p /apps/jvm-logs
ADD ./minioUploadFile ${HOME}/minioUploadFile
ADD ./config.toml ${HOME}/config.toml
ADD ./upload.sh ${HOME}/upload.sh
ADD target/oom-test-*.jar ${HOME}/app.jar
WORKDIR ${HOME}

EXPOSE 8080

CMD [ "sh", "-c", "java $JAVA_OPTS $OOM_OPS -Djava.security.egd=file:/dev/./urandom $JAVA_AGENT_OPTS -jar app.jar" ]

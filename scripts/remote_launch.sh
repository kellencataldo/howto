#!/bin/bash


# parameter 1 is the compiler serverht binary
# parameter 2 is the client certificate

set -e

test -e $1 || echo "serverht file not found"
test -e $2 || echo "client certificate not found"


ssh -i $SECRETS/aws-instance.pem -o StrictHostKeyChecking=no $AWS_INSTANCE "ps -A | grep serverht | awk '{print \$1}' | xargs kill -2"

scp -i $SECRETS/aws-instance.pem $1 $AWS_INSTANCE:/home/ec2-user/serverht/serverht
scp -i $SECRETS/aws-instance.pem $2 $AWS_INSTANCE:/home/ec2-user/serverht/client.crt

ssh -i $SECRETS/aws-instance.pem -o StrictHostKeyChecking=no $AWS_INSTANCE "nohup /home/ec2-user/serverht/serverht \
    -clientcert /home/ec2-user/serverht/client.crt \
    -serverkey /home/ec2-user/htserer/server.key \
    -servercert /home/ec2-user/serverht/server.crt \
    -port 8080 \
    -logfile /home/ec2-user/serverht/logs/serverht.log 1>/home/ec2-user/serverht/logs/nohup.log 2>/home/ec2-user/serverht/logs/nohup.log &"

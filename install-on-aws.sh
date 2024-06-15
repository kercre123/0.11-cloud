#!/bin/bash

sudo ./deploy-create.sh arm64

scp -i $1 deploy.tar.gz admin@$2:/wire/

ssh -i $1 admin@$2 "cd /wire && sudo rm -rf new-deploy && mkdir new-deploy && mv deploy.tar.gz new-deploy/ && cd new-deploy && tar -zxvf deploy.tar.gz && ./install.sh -u /wire/011-cloud && cd .. && rm -rf new-deploy && echo success"

#!/bin/sh

 if [ "$INITIALIZE" == "true" ]
 then
    /clair -config /config/config.yaml
else
    /clair -config /config/config.yaml & 
    echo "security-scanner" && docker login $REGISTRYSERVER -u $REGISTRYUSERNAME -p $REGISTRYPASSWORD 
    docker pull $REGISTRYIMAGE
    cd /clairctl
    mkdir -p /clairctl/reports/html && chmod -R 777 /clairctl/reports/html
    clairctl --config=clairctl.yml health 
    clairctl --config=clairctl.yml analyze -l $REGISTRYIMAGE  
    clairctl --config=clairctl.yml report -l $REGISTRYIMAGE 
fi

sleep 6000

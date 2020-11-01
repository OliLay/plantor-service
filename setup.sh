#!/bin/bash
# create docker volumes
mkdir /srv/plantor
mkdir /srv/plantor/grafana
mkdir /srv/plantor/grafana/plugins
mkdir /srv/plantor/grafana/data

# chown the grafana user
sudo chown -R 472:472 /srv/plantor/grafana

echo "Setup successful!"
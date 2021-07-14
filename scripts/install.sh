sudo add-apt-repository ppa:maxmind/ppa
sudo apt install libmaxminddb0 libmaxminddb-dev mmdb-bin
cd scripts
bash get_latest_exit_nodes.sh
bash get_cloud_host.sh
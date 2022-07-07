#! /bin/bash

echo '\033[0;35m
  ██████  █    ██   ██████▓██   ██▓    ███▄    █  ▒█████  ▓█████▄ ▓█████ 
▒██    ▒  ██  ▓██▒▒██    ▒ ▒██  ██▒    ██ ▀█   █ ▒██▒  ██▒▒██▀ ██▌▓█   ▀ 
░ ▓██▄   ▓██  ▒██░░ ▓██▄    ▒██ ██░   ▓██  ▀█ ██▒▒██░  ██▒░██   █▌▒███   
  ▒   ██▒▓▓█  ░██░  ▒   ██▒ ░ ▐██▓░   ▓██▒  ▐▌██▒▒██   ██░░▓█▄   ▌▒▓█  ▄ 
▒██████▒▒▒▒█████▓ ▒██████▒▒ ░ ██▒▓░   ▒██░   ▓██░░ ████▓▒░░▒████▓ ░▒████▒
▒ ▒▓▒ ▒ ░░▒▓▒ ▒ ▒ ▒ ▒▓▒ ▒ ░  ██▒▒▒    ░ ▒░   ▒ ▒ ░ ▒░▒░▒░  ▒▒▓  ▒ ░░ ▒░ ░
░ ░▒  ░ ░░░▒░ ░ ░ ░ ░▒  ░ ░▓██ ░▒░    ░ ░░   ░ ▒░  ░ ▒ ▒░  ░ ▒  ▒  ░ ░  ░
░  ░  ░   ░░░ ░ ░ ░  ░  ░  ▒ ▒ ░░        ░   ░ ░ ░ ░ ░ ▒   ░ ░  ░    ░   
      ░     ░           ░  ░ ░                 ░     ░ ░     ░       ░  ░
                           ░ ░                             ░'

echo '\033[0m '
echo "Welcome to Susy Chainlink Auditor installation guide!"

while true; do
    read -p  "Do you want to start installation process?  [Y/n] " yn 
    case $yn in
        [Yy]* ) apt-get update; break;;
        [Nn]* ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
done


docker -v

if [ $? -eq 0 ]; then
  echo '\033[0;33mdocker already installed'
  echo '\033[0m '
else 
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
  echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
  sudo apt-get update
  sudo apt-get install docker-ce docker-ce-cli containerd.io 
  docker -v
  if [ $? -eq 0 ]; then 
    echo '\033[0;33mdocker successfully installed'
    echo '\033[0m '
    else
    echo '\033[0;31mdocker installation aborted. Please check logs and install it manualy'
    echo '\033[0m '
    exit
  fi
fi

docker-compose -v

if [ $? -eq 0 ]; then
  echo '\033[0;34mdocker-compose already installed'
  echo '\033[0m '
else 
  sudo apt-get install docker-compose
  docker-compose -v
  if [ $? -eq 0 ]; then 
    echo '\033[0;32mdocker-compose successfully installed'
    echo '\033[0m '
    else
    echo '\033[0;31mdocker-compose installation aborted. Please check logs and install it manualy'
    echo '\033[0m '
    exit
  fi
fi

git clone https://github.com/SuSy-One/susy-chainlink-auditor
echo '\033[0m '

cd susy-chainlink-auditor/
mkdir chainlink-1
cd chainlink-1


touch .password
echo "Create new inner password(A password requires to have at least 3 numbers and 3 symbols):"
read passw
echo $passw > .password

echo '\033[0m '

touch .api
echo "Enter your email:"
read mail 
echo $mail > .api

echo '\033[0m '

echo "Create your password:"
read password 
echo $password >> .api

  echo -e '\033[0;34mNode launch...' 
  echo '\033[0m '


cd /var/www/susy-chainlink-auditor && docker-compose up -d

cd /var/www/susy-chainlink-auditor && docker-compose ps

echo '\033[0m '

  echo '\033[0;33mUse http://localhost:7788 to access Chainlink Operator GUI if you are running the node on a local machine. If not - use ssh tunneling to be able to launch console using web browser: ssh -L 7788:localhost:7788 root@your.server.ip.address'
  echo '\033[0m '

echo 'Your node info:'

sudo docker-compose exec node1 chainlink admin login

echo '\033[0m '

sudo docker-compose exec node1 chainlink keys eth list

echo '\033[0m '

sudo docker-compose exec node1 chainlink keys ocr list

echo '\033[0m '

sudo docker-compose exec node1 chainlink keys p2p list

echo '\033[0m '


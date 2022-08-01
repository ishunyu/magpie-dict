# Setup Nginx as reverse proxy
[DigitalOcean - How To Deploy a Go Web Application Using Nginx on Ubuntu 18.04](https://www.digitalocean.com/community/tutorials/how-to-deploy-a-go-web-application-using-nginx-on-ubuntu-18-04)

1. Install Nginx
```
sudo apt update
sudo apt install nginx
```
2. Adjust firewall (if needed)
```
sudo ufw app list
sudo ufw allow 'Nginx HTTP'
sudo ufw status
```
3. Checking your Web Server

```
systemctl status nginx
```
4. Setting Up a Reverse Proxy with Nginx
```
cd /etc/nginx/sites-available
```
```
sudo nano your_domain
```
Add the following into file
```
server {
    server_name your_domain www.your_domain;

    location / {
        proxy_pass http://localhost:9990;
    }
}
```
```
sudo ln -s /etc/nginx/sites-available/your_domain /etc/nginx/sites-enabled/your_domain
```
```
sudo nginx -s reload
```
5. Forward client IP
```
proxy_set_header   Host               $host;
proxy_set_header   X-Real-IP          $remote_addr;
proxy_set_header   X-Forwarded-Proto  $scheme;
proxy_set_header   X-Forwarded-For    $proxy_add_x_forwarded_for;
```
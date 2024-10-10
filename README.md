<h1> Initial configuration </h1>

First things first, you'd need to clone the main.go file and switch into the downloaded folder.

```
> git clone https://git.voltagexd.gay/voltage/Go-ShareX-Server
> cd Go-ShareX-Server
```

Then, you can begin to configure the script or adjust the code to your liking.

```
> nano main.go
```

Basically, the only thing you'd want to change is the port or name of the uploads directory. However, you can just stick to the initial configuration.

<h1> Configuring nginx </h1>

It's only useful if you are using a domain for your server. Create a new nginx configuration file.

```
> nano /etc/nginx/sites-available/hosts  # change the name to your liking
```

Then, you can just copy and paste my configuration.

```
server {
    listen 80;
    server_name domain.here;

    location /upload {
        proxy_pass http://localhost:8040/upload;  # change the port
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Then, create a symlink of this file in sites-enabled directory

```
> sudo ln -s /etc/nginx/sites-available/sharex /etc/nginx/sites-enabled/
> sudo nginx -t
> sudo systemctl restart nginx
```

If you want to use ssl:

```
> apt install certbot
> apt install python3-certbot-nginx  # make sure the certbot and nginx module are installed
> certbot --nginx -d yourdomain.here
```

Certbot will automatically setup all the ssl certificates for your domain.

Finally, restart nginx

```
> systemctl restart nginx
```

<h1> Make it a service </h1>

By creating sharex-srv.service file in /etc/systemd/system/ directory

```
> nano /etc/systemd/system/sharex-srv.service
```

And pasting the following config:

```
[Unit]
Description=ShareX File Upload Server
After=network.target

[Service]
ExecStart=/path/to/your/sharex-server-binary
WorkingDirectory=/path/to/your/sharex-server
Restart=always
User=yourusername
Group=yourgroupname
Environment=PORT=8040  # change the port to the one you have in your main.go file

[Install]
WantedBy=multi-user.target
```

Finally, build the server binary by running the following command (make sure Go is installed)

```
> go build main.go
```

And run your server by running

```
> systemctl start sharex-srv
> systemctl enable sharex-srv
```

<h1> ShareX configuration </h1>

Create a new custom uploder with following settings:

```
Name: Any
Destination type: Image uploader
Method: POST
Request URL: youripordomain.com/upload
URL Parameters: None
Body: Form data
Headers: {
    Name: X-File-Name, Value: {filename}
}
File from name: file
URL: {response}
```

That's it, now you have your very own ShareX host! Go makes it fast, ShareX makes it useful.
linx-client
======

Simple client for [linx-server](https://github.com/andreimarcu/linx-server) 

Uses a json configuration and logfile (for storing deletion keys of uploads)   

### Features  

- Upload file and store deletion key  
- Upload from stdin and store deletion key  
- Overwrite file using stored deletion key  
- Delete file using stored deletion key  
- Sha256sum client-server matching  


Get release and run
-------------------
1. Grab the latest binary from the [releases](https://github.com/andreimarcu/linx-client/releases)
2. Run ```./linx-client...```


Configuration
-------------

When you first run linx-client, you will be prompted to configure the instance url, logfile path and api keys. 

```
$ ./linx-client
Configuring linx-client

Site url (ex: https://linx.example.com/): https://linx.example.com/
Logfile path (ex: ~/.linxlog): ~/.linxlog
API key retreival command (leave blank for plain token or if instance is public, ex: pass show linx-client): pass show 'linx-client'
Configuration written at ~/.config/linx-client.conf
```


### With API key

Your API key can either be stored directly in plain text in your
`~/.config/linx-client.conf` file, such as:

```json
{
    "apikey": "your-api-key-here",
    "logfile": "/home/you/.linxlog",
    "siteurl": "https://linx.example.com/"
}
```

Or you can use an external tool, such as [pass](https://www.passwordstore.org/)
to store it as encrypted.

If you specify the config key `apikeycmd`, it will take precedence over the
`apikey` config key, and will run and use the first line of output from that
command when it performs any request.

```json
{
    "apikey": "",
    "apikeycmd": "pass show linx-client",
    "logfile": "/home/you/.linxlog",
    "siteurl": "https://linx.example.com/"
}
```

The `apikeycmd` config key does support quoting. Specifying
`pass show 'my key'` will resolve in running command `pass` with the 2
arguments, `show`, and `my key`. Whereas `pass show my key` would resolve as
the command `pass` and the 3 arguments `show`, `my`, and `key`. Different
programs handle argument grouping differently. But for [pass](https://www.passwordstore.org/)
in particular the latter would be an invalid command as it only expects 1 pass
name.


### Without API key (public instance)

By leaving both `apikey` and `apikeycmd` empty, linx-client will assume the
target instance is publicly writable and will then perform its requests without
an API key.


Usage
----- 

#### Upload file(s)

```
$ linx-client path/to/file.ext
https://linx.example.com/file.ext
```

Options  

- ```-f file.ext``` -- Specify a desired filename (if different from the actual one)  
- ```-r``` -- Randomize filename  
- ```-e 600``` -- Time until file expires in seconds  
- ```-deletekey mysecret``` -- Specify deletion key
- ```-o``` -- Overwrite file if you have its deletion key
- ```-accesskey mykey``` -- Specify access key
- ```-c myconfig.json``` -- Use non-default config file (can be useful if using more than one linx-server instance). This option will create a config if file does not exist.
- ```-no-cb``` -- Disable automatic insertion into clipboard
- ```-selif``` -- Return selif url (direct url)

#### Upload from stdin
```
$ echo "hello there" | linx-client -  
https://linx.example.com/random.txt  
```  

Note: you can specify the ```-f``` flag to specify a filename as such:  

```
$ echo "hello there" | linx-client -f hello.txt -  
https://linx.example.com/hello.txt  
```  


#### Overwrite file
Assuming you have previously uploaded ```file.ext``` using linx-client (so that you have its deletion key), you can replace the file as such:

```
$ linx-client -o file.ext  
https://linx.example.com/file.ext  
```  

#### Delete file(s)

```
$ linx-client -d https://linx.example.com/file.ext  
Deleted https://linx.example.com/file.ext  
```

License
-------
Copyright (C) 2015 Andrei Marcu

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

Author
-------
Andrei Marcu, http://andreim.net/

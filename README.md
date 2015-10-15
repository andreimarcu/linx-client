linx-client
======

Simple client for [linx-server](https://github.com/andreimarcu/linx-server) 

Uses a json configuration and logfile (for storing deletion keys of uploads)   


Configuration
-------------

When you first run linx-client, you will be prompted to configure the instance url, logfile path and api keys. 

```
$ ./linx-client  
Configuring linx-client  
  
Site url (ex: https://linx.example.com/): https://linx.example.com/  
Logfile path (ex: ~/.linxlog): ~/.linxlog  
API key (leave blank if instance is public):  
  
Configuration written at ~/.config/linx-client.conf  
```

Usage
----- 

#### Upload file(s)

```
$ linx-client path/to/file.ext
https://linx.example.com/file.ext
```

Options  

- ```-r``` -- Randomize filename  
- ```-e 600``` -- Time until file expires in seconds  
- ```-deletekey mysecret``` -- Specify deletion key


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

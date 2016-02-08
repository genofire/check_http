# check_http

## Example .checkhttprc

```
---
ipv4: true
ipv6: true
domains:
- domain: google.de
  regex: "<title>.*</title>"

```
## Output
0 = regex is okay
W = request is okay, but not regex
30 = a http 30* recieved
E = everything else

```
+-----------------------------+------+------+----------+----------+
|           DOMAIN            | IPV4 | IPV6 | SSL-IPV4 | SSL-IPV6 |
+-----------------------------+------+------+----------+----------+
| google.com                  | E    | E    | E        | E        |
+-----------------------------+------+------+----------+----------+
```

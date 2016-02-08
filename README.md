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

| Code | Result |
|:----:|--------|
|    0 | regex is okay |
|   30 | a http 30* recieved |
| W    | request is okay, but not regex |
| E    | everything else "error" |


```
+-----------------------------+------+------+----------+----------+
|           DOMAIN            | IPV4 | IPV6 | SSL-IPV4 | SSL-IPV6 |
+-----------------------------+------+------+----------+----------+
| google.com                  | E    | E    | E        | E        |
+-----------------------------+------+------+----------+----------+
```

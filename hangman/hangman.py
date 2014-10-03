#!/usr/bin/env python3.4

import sys

try:
    import json, requests
except ImportError as e:
    print(e)
    sys.exit(1)

url = 'http://richmond.cookgetsitdone.com:8080/hangman'

hdrs = { 'content-type': 'application/json'}

payload1 = {"Cmd": "NEW"}
payload2 = {
    "Cmd": "P1T", 
    "Play": "i", 
    "Gid": 1, 
    "Auth": ''
}

try:
    resp = requests.post(url, data=json.dumps(payload1))
except requests.packages.urllib3.exceptions.ProtocolError as e:
    print(e)
    sys.exit(2)

try:
    if isinstance(sys.argv[1], str):
        sys.exit(0)
except IndexError:
    pass

payload2["Gid"] = resp.json()["Game"]
payload2["Auth"] = resp.json()["Cred"]

print(resp.json()["Hint"])

wrong = []
doit = lambda x: chr(x) if x>0 else '-'
while len(wrong) < 7:
    print(''.join([doit(l) for l in resp.json()["Curr"]]))
    payload2["Play"] = input("Letter: ")
    if not payload2["Play"]:
        break
    try:
        resp = requests.post(url, data=json.dumps(payload2))
    except requests.packages.urllib3.exceptions.ProtocolError as e:
        print(e)
        break
    try:
        wrong = resp.json()["Missed"]
    except KeyError as e:
        print("Oops")

import time
import hmac
import hashlib
import requests

secret = ""
url = ""

timestamp = str(int(time.time()))

sign = hmac.new(
    secret.encode('utf-8'),
    timestamp.encode('utf-8'),
    hashlib.sha256
).hexdigest()

headers = {
    "X-Timestamp": timestamp,
    "X-Signature": sign
}

response = requests.get(url, headers=headers)
print(response.text)
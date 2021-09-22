from Crypto.PublicKey import RSA
from Crypto.Cipher import PKCS1_v1_5 as Cipher_PKCS1_v1_5
from base64 import b64decode,b64encode
pubkey = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCkq3XbDI1s8Lu7SpUBP+bqOs/MC6PKWz6n/0UkqTiOZqKqaoZClI3BUDTrSIJsrN1Qx7ivBzsaAYfsB0CygSSWay4iyUcnMVEDrNVOJwtWvHxpyWJC5RfKBrweW9b8klFa/CfKRtkK730apy0Kxjg+7fF0tB4O3Ic9Gxuv4pFkbQIDAQAB"
msg = "Z2PZJzHWs97My96VEJSTQE5b4q66blOq"
keyDER = b64decode(pubkey)
keyPub = RSA.importKey(keyDER)
cipher = Cipher_PKCS1_v1_5.new(keyPub)
cipher_text = cipher.encrypt(msg.encode())
emsg = b64encode(cipher_text)
print(emsg)
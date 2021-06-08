## Get started

Run application

```
make up
```

Note: Application will use 8080 port

---

## Auto tests

```
make test
```

---

## Manual test

1. Get mnemonic

```
GET 'localhost:8080/api/v1/btc/wallet/mnemonic'
```

2. Create HD SegWit Address

```
POST 'localhost:8080/api/v1/btc/wallet/hd/segwit'

Headers:
{
    "Content-Type": "application/json"
}

BOdy:
{
    "seed": [bytes...],
    "path": (string)
}

Example body:
{
    "seed": [244,184,4,62,59,59,77,11,158,60,124,218,129,214,134,140,51,26,174,204,128,85,93,199,178,208,237,206,107,115,234,80,169,29,103,88,111,116,97,205,70,202,204,238,110,36,10,89,138,154,170,48,99,205,217,190,198,90,61,36,211,170,85,27],
    "path": "m/84'/0'/0'/0/0"
}
```

3. Create MUltiSig P2KH Adress

```
POST 'localhost:8080/api/v1/btc/wallet/multisig'

Headers:
{
    "Content-Type": "application/json"
}

BOdy:
{
    "n": (int),
    "m": (int),
    "public_keys": [string...]
}

Example body:
{
    "n": 3,
    "m": 2,
    "public_keys": [
        "04a882d414e478039cd5b52a92ffb13dd5e6bd4515497439dffd691a0f12af9575fa349b5694ed3155b136f09e63975a1700c9f4d4df849323dac06cf3bd6458cd","046ce31db9bdd543e72fe3039a1f1c047dab87037c36a669ff90e28da1848f640de68c2fe913d363a51154a0c62d7adea1b822d05035077418267b1a1379790187","0411ffd36c70776538d079fbae117dc38effafb33304af83ce4894589747aee1ef992f63280567f52f5ba870678b4ab4ff6c8ea600bd217870a8b4f1f09f3a8e83"
    ]
}
```

---

### Library used
- github.com/gorilla
- github.com/btcsuite
- github.com/tyler-smith/go-bip32
- github.com/tyler-smith/go-bip39
- golang.org/x/crypto
- github.com/stretchr/testify


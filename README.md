# Celestia DA Proxy API

Descrption...

## Example

1. Submit Image Blob

```bash

 curl --location 'http://localhost:3000/submit_formdata' \
--form 'image=@"/Users/jeongseup/Downloads/CelestiaDragonsNFT/DALL·E 2024-05-15 19.11.55 - A cute, animated-style dragon designed to be stored on a blockchain named Celestia. The dragon has big, sparkling eyes and a small, friendly smile. It.webp"'

#{"success":true,"result":{"hash":"97110C95BB4FD3BDD7520744A2FA8C03DC034DDE2C668647784F32A4BAE14E8C","submitted_height":1844378}}%
# curl -X GET 'http://localhost:3000/97110C95BB4FD3BDD7520744A2FA8C03DC034DDE2C668647784F32A4BAE14E8C'
```

2. Submit Metadata Blob With Namespace Key

```bash
curl -sS -X POST http://localhost:3000/submit_metadata -H "Content-Type: application/json" -d '{
  "namespace_key": "CelestiaDragonsMetaData",
  "metadata": {
    "description": "Celestia DA Based Dragon NFT Collection",
    "image": "https://nftinfo.online/97110C95BB4FD3BDD7520744A2FA8C03DC034DDE2C668647784F32A4BAE14E8C",
    "name": "Celestia First DA Dragon",
    "attributes": [
      {
        "trait_type": "Color Palette",
        "value": "Pastel pinks, blues, and purples"
      },
      {
        "trait_type": "Environment",
        "value": "Clouds with twinkling stars"
      },
      {
        "trait_type": "Disposition",
        "value": "Friendly and playful"
      },
      {
        "trait_type": "Special Feature",
        "value": "Translucent, ethereal wings"
      },
      {
        "trait_type": "Magic Power",
        "value": "Can manipulate weather patterns"
      }
    ]
  }
}' | jq .
```

3. Retrieve Metadata By Namespace Key & Index

```bash

curl -sS -X GET http://localhost:3000/Q2VsZXN0aW/1
```

```bash

# retrieve blobs
curl -X POST http://localhost:3000/retrieve_blob -H "Content-Type: application/json" -d '{"retrieve_height":"1830482"}'

baseuri = nftinfo.online/jdfklasljfsflkasj -> call <- response was?


curl -X POST http://hackathemy.me:3000/submit_blob -H "Content-Type: application/json" -d '{"data":"SGVsbG8gV29ybGQ="}'
curl -X POST http://hackathemy.me:3000/retrieve_blob -H "Content-Type: application/json" -d '{"retrieve_height":"1830482"}'

# 1. baseuri = nftinfo.online/metadata/jdfklasljfsflkasj -> call <- response was?
# 1. post image data to da
curl --location 'http://localhost:3000/submit_formdata' \
--form 'image=@"/Users/jeongseup/Downloads/CelestiaDragonsNFT/DALL·E 2024-05-15 19.11.55 - A cute, animated-style dragon designed to be stored on a blockchain named Celestia. The dragon has big, sparkling eyes and a small, friendly smile. It.webp"'

# 2. check posted image data from da
curl -X GET 'http://localhost:3000/retrieve_blob?height=1834058' -H "Content-Type: application/json"

# 3. post metadata to da
curl -X POST http://localhost:3000/submit_metadata -H "Content-Type: application/json" -d '{
  "metadata": {
    "description": "Celestia DA Based Dragon NFT Collection",
    "external_url": "https://openseacreatures.io/3",
    "image": "https://nftinfo.online/97110C95BB4FD3BDD7520744A2FA8C03DC034DDE2C668647784F32A4BAE14E8C",
    "name": "Celestia First DA Dragon",
    "attributes": [
      {
        "trait_type": "Color Palette",
        "value": "Pastel pinks, blues, and purples"
      },
      {
        "trait_type": "Environment",
        "value": "Clouds with twinkling stars"
      },
      {
        "trait_type": "Disposition",
        "value": "Friendly and playful"
      },
      {
        "trait_type": "Special Feature",
        "value": "Translucent, ethereal wings"
      },
      {
        "trait_type": "Magic Power",
        "value": "Can manipulate weather patterns"
      }
    ]
  }
}'

curl -X GET 'http://localhost:3000/retrieve_blob?height=1834153' -H "Content-Type: application/json"
```

## Test

```bash
# test for receving form-data
curl --location 'http://localhost:3000/test_receive_formdata' \
--form 'image=@"/Users/jeongseup/Downloads/CelestiaDragonsNFT/DALL·E 2024-05-15 19.11.55 - A cute, animated-style dragon designed to be stored on a blockchain named Celestia. The dragon has big, sparkling eyes and a small, friendly smile. It.webp"'

# test for receiving json-data
curl -sS -X POST http://localhost:3000/test_receive_jsondata -H "Content-Type: application/json" -d '{
  "namespace_key": "CelestiaDragonsMetaData",
  "metadata": {
    "description": "Celestia DA Based Dragon NFT Collection",
    "external_url": "https://openseacreatures.io/3",
    "image": "https://nftinfo.online/images/first.png",
    "name": "Celestia First DA Dragon",
    "attributes": [
      {
        "trait_type": "Color Palette",
        "value": "Pastel pinks, blues, and purples"
      },
      {
        "trait_type": "Environment",
        "value": "Clouds with twinkling stars"
      },
      {
        "trait_type": "Disposition",
        "value": "Friendly and playful"
      },
      {
        "trait_type": "Special Feature",
        "value": "Translucent, ethereal wings"
      },
      {
        "trait_type": "Magic Power",
        "value": "Can manipulate weather patterns"
      }
    ]
  }
}' | jq .

```

curl -X GET http://hackathemy.me/health

## Depedency

```bash
brew install potrace
```

### References

https://jaredwinick.github.io/base64-image-viewer/?ref=tiny-helpers

https://www.based64.xyz/

https://docs.celestia.org/developers/golang-client-tutorial

https://docs.opensea.io/docs/metadata-standards

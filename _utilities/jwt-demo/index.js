var jwt = require('jsonwebtoken');
var privateKey = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDAhCwsYpPX96wH
mOT9u/a/I6J6HteMeymXKyFEEzS7hj9IOmkFvJMXMamO0vmTxntW9EWVIRAhDP+f
9A1UirEGhFe6cGhho+NVEY8uTc1jEsMpHQThlVlluMJ6aMpXeaqC+bCh5xz10iVv
EHsmqG0ajQVP0lgkTwbA8J228Lj3jd22hcj3reh9JzMBGDfKXwYuPXdRCBYWUPEd
UPRBrvMFQhmY0plBwm8hsuAgG9XpubkWG/0nEZA2eqohLtXcGcr8+yZTYa9U6RrU
I6CAJt86ZoIY+kVoHQ9rKJxHbaoxEV68PVgrogTel66vwubQZPYVtAAAta8M7Ahj
y6OaaP3fAgMBAAECggEAAQu2wbb3XqD1ykTGWap/uKIU69znBthUbYHKeHgyPSKJ
jXbjwyg3FFUsup3ixS661MoW5qx7TfKoINJOkDsAoyxs3ZJmvsEJJxklUfcSOouL
i9mggSkyhx7tz4RqNPdwOa3pi7eZrKwrj+LjCF652P8THE/kMq5LNxkXgjrkhlO4
iyVQhNfv6cGWljIdWC61UJu0KNX+SOFsexveMaRiWv525LI3svCv+ONJbFPRV4DN
LoDAqMFHoOSrNEC7WRHYuKT+mnMShtoXUfbR5W20Mgca0JJLbH3UfmCX6n93Z+xY
SDNJWijSZInbJ7pbonThcUs5PbC42W6h+xv9TdnNYQKBgQDfhYkA61npy+640yTp
2xr2LsKuSNB3SJyL9tbeKGYvvqwLlKX9pgqeBDz4+Il8XVPDjGMjVCoslcF81CyA
vLuozq1fYWwp8vNAVmnc12J7LHuskZ4HcCkBQkpc7ynnhR17JkaX7WFs8GzLeprE
vkfd/xhvUaI24A3E8HZwocJA6wKBgQDcfVBwNEURnAGUOPOiky7Z1no7s0DhxUj9
N3wGoX9qjegdAH+9SuafFV062QQ0pL6in1en10MoiX+UMEqHmjGzFg8JgJMFb+Xo
1VrxgMj50gQdq/Sx5csTw2mKyON4UvaeKPiCzIJajDDcUNUs/Zu6es2EhWUPE1pe
9jsDoCAZ3QKBgQDa2V0N6GFtlz2R+zHOQrgASEJW8HYkBJU6OSGh/L4oizDaWd24
HuWQV6f3QSkj+iC0evTqN1LfunTqqrc0CRZYLpvzomiMHhLrcHBQSRZkcWZZzW2D
7N5JHEXA/m3yABSgahZ+VG6qgjCTfeShM4kcI9Mh0zTXM8Jni+T7XHXcpwKBgQC2
L7WvgQExM776QrTNuOAVj2sguVT7OJC+6oHI2Nj3qpoInMjwGFvHR1fpsDgRZ689
oHxFa1FKxZJtWBm9QmOenrN+HoddDsDiSqkCtG9cPXS5L8TY2g+bHPSwgJ20Zpjw
xtnQ+jsbposZAJGkw0lSJPZ8cdy3QD6ECOFqdX0Q4QKBgQDFxLkKMlEgaat1+M5J
YMf/Z2Jy8EHhPJazBoTJ0qf1/IQFW7Mil1urP7O1HttHoOwf+6VTg81/5Q+IahOZ
pZmmZCGIyEAS8CmzVgvxEWoVYxkMOU1rFKI96HqiUbP1NsieT/47sbAEtxtWxiRX
SlS1xnTYLjU0paQm/JIYshr/ww==
-----END PRIVATE KEY-----`
var user_id = '3b2774bf-bbd5-4d65-a3f4-303d9e881466';
var token = jwt.sign(
    {
        "sub": user_id,
        "aud": "lu2k7ptfjb4gh05nkg1kd3u9v",
        "cognito:groups": ["Doctors"],
        "iss": "https://cognito-idp.ap-southeast-2.amazonaws.com/ap-southeast-2_gfSuuHw6e",
        "version": 2,
        "token_use": "access",
        "auth_time": (new Date().getTime())/1000, //"issued_at_timestamp_in_seconds",
        "exp": (new Date().getTime())/1000 + 3600, //"expires_at_timestamp_in_seconds",
        "iat": (new Date().getTime())/1000, //"issued_at_timestamp_in_seconds",s
        "jti": user_id,
        "username": "voanh85@gmail.com"
    },
    privateKey,
    {
        algorithm: 'RS256',
        keyid: "TJAVF4GH",
    }
)

console.log(token)
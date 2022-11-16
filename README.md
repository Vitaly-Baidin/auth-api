# Start/stop server
**Start:**
```bash
make start
```
**Stop:**
```bash
make stop
```

# Endpoints

- localhost:8080/v1/auth/register register
```json
{
	"login": "vitaly",
	"email": "sdaw@mail.com",
	"phone": 89004002000,
	"password": "12345678"
}
```
- localhost:8080/login login
```json
{
	"email": "sdaw@mail.com",
	"password": "12345678"
}
```
- localhost:8080/v1/ping (need add header param "Bearer-Token":"{Token access from login req}") test token
- localhost:8080/v1/auth/refresh refresh access login
```json
{
	"token": "{refresh token from login req}"
}
```
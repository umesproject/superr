# Superr

Superr is a structured error library for golang.

It was inspirated by: https://middlemost.com/failure-is-your-domain/

### Getting started
Add 
`go get github.com/umesproject/superr`

### Example 1

```go
func Login() error {

    // good pratice: define the operation name 
	// at the start of the function
    const op superr.Op = "user.Login"
	

    err := Validate()
	if err != nil {
		return superr.E(err, op)
	}

	return nil
}
```

### Example 2  ( with extra fields passed to error )

```go
func Login() error {

    // good pratice: define the operation name 
	// at the start of the function
    const op superr.Op = "user.Login"
	

    err := Validate()
	if err != nil {
		return superr.E(err, op)
	}

	return nil
}

func Validate() error {
	const op superr.Op = "utils.Validate"
    // you can pass extra fields in the superr.Fields struct
	return superr.E(op, superr.Fields{"username": "test", "usedRegex": false})
}

func main(){
	err := Login()

	if err != nil {
		superr.Log(err)
	}
}
```
Output examplpe 2:
```json
{
  "caller": "File: main.go  Function: main.Login Line: 17",
  "extraFields": {
    "ok": true,
    "prova": 12
  },
  "level": "info",
  "msg": "",
  "stackTrace": [
    "user.Login",
    "utils.Validate"
  ],
  "time": "2021-02-03T18:32:42+01:00"
}
```

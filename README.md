Form3 Take Home Exercise
----

Overview
----
This is a client library in Go to access Form3 fake account API. Please note I'm new to the golang programming language 
although I've used it in a couple of small personal projects.

The library has two "layers", one ('data' package) contains the logic to access the external API and the other the logic 
to build the data which ultimately will be received by any client that uses this library. 
I think is good and simple way of separating concerns, keeping the general implementation simple, concise, 
readable and testable.

The **data** package connects to the Account API via an environment variable (ACCOUNT_API_ADDR). I thought it would be 
good to not keep it hardcoded as is something we may want to modify without requiring to build the code just for that. 
I've decided to use a 'dto' due to the API be something this project does not control. This way the type 'AccountDto'
is used to help transferring data between the API and our library but not to send it to clients as a final result.
Another point related to this independence among layers is the interface _'AccountApiGateway'_ which states a contract 
that can be used without any fear of breaking the top layer's implementation even if the _data_ package needs to change. 
This package does not make use of any cache mechanism as I considered it to be an advanced feature and therefore not 
required for now. This means for each call there will be an http request (given all the input is correct).

The **form3** package represent a higher level layer and tries to handle the input given by any 
client and resolve it into the final result. This library represents the concept of an 
account by the type 'Account'. The decision of returning a pointer in the _create_ and _get_ action functions is based
upon:

* Mutability - Client code can modify our returned data throughout their code.
* Signify true absence - If no account is to be returned the library makes use of the _zero value_ for pointers (nil).
* Large structures - Account type has already quite some fields and could be even bigger. An address to this data would 
  make it more efficient than the process of copying values (if this library used values instead of pointers).
  
I've tried to follow as many idiomatic concepts as possible while using Go in this project. One of those example can be 
seen in the functions that help to create a new instance of some struct _(NewXxxx)_. This library makes use of two 
external packages, the testing framework [https://github.com/matryer/is](is) which in my opinion makes the tests more 
readable and [github.com/google/uuid](uuid) to manage the UUID types.

How to use
----

### Run tests:

```
docker-compose up
````

### Usage:

```go
func CreateAccount(info *Account) (*Account, error)
```
Create a new account with the given info. Returns an error if a problem occurs while trying to create the new account.

```go
func DeleteAccount(id string, vrs int) error
```

Delete the account with a given id. Id must be a valid uuid type.
Returns an error if a problem occurs while trying to delete the account.

```go
func GetAccount(id string) (*Account, error)
```
Retrieves an account by the given id. Id must be a valid uuid type.
Returns an error if a problem occurs while trying to delete the account with the given id.



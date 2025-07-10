## Go JSON Polymorphism

Go lacks native support for sum types (unlike Rust's enums with exhaustive matching), forcing developers to choose between workarounds when modeling polymorphic data structures. This experiment compares two approaches to handling polymorphic JSON structures in Go.

### [Sum Type Approach](./sum.go)

- Uses interfaces and a registry pattern
- Provides strong type safety through distinct credential types
- Requires custom JSON unmarshaling logic
- More complex but extensible architecture

```console
$ go run . -sum
=== Using Sum Type Approach ===
USERNAME-PASSWORD
Unmarshalled: {CredentialsType:username-password CredentialsBody:Username:john Password:secret123}
Marshalled: {"credentials-type":"username-password","credentials-body":{"username":"john","password":"secret123"}}

TOKEN
Unmarshalled: {CredentialsType:token CredentialsBody:Token:abc123xyz}
Marshalled: {"credentials-type":"token","credentials-body":{"token":"abc123xyz"}}
```

_This is just an example, you shouldn't print passwords willy-nilly on a production system._

### [Flat Struct Approach](./flat.go)

- Uses a single struct with conditional validation
- Uses `go-playground/validator` tags for validation
- Simpler implementation with built-in JSON support
- All fields visible but only relevant ones populated

```console
$ go run . -flat
=== Using Flat Struct Approach ===
USERNAME-PASSWORD
Unmarshalled: {Credentials:{Type:username-password Username:alice Password:mypassword Token:}}
Marshalled: {"credentials":{"type":"username-password","username":"alice","password":"mypassword"}}

TOKEN
Unmarshalled: {Credentials:{Type:token Username: Password: Token:xyz789abc}}
Marshalled: {"credentials":{"type":"token","token":"xyz789abc"}}
```

_This is just an example, you shouldn't print passwords willy-nilly on a production system._

### Differences

- **Complexity**: The sum type approach requires interfaces, constructors, registry patterns, and custom unmarshaling. The flat approach uses only validation tags.
- **Type Safety**: Sum types provide compile-time guarantees about which types exist, but field access still requires runtime type assertions. Flat structs rely on runtime validation but offer direct field access.
- **Extensibility**: Sum types require compile-time registration of new types but provide cleaner separation. Flat structs need additional fields and validation rules but can handle unknown types more gracefully.
- **Performance**: Sum types incur costs from reflection during unmarshaling and interface calls. Flat structs avoid interface dispatch and reduce allocations.

### When to Use Each Approach

- **Sum types**: When you need strict type boundaries, have a stable set of credential types, or are building a library where type safety is paramount.
- **Flat structs**: When you prioritize simplicity, need to handle evolving schemas, or are building applications where performance and maintainability matter more than type purity.

# somaadm logout

This command can be used by a user to revoke active authentication
tokens.

By default it revokes the token currently used by the somaadm client. If the all flag is provided, it revokes all tokens for the account that were issued prior to the revocation.

If the client currently does not have an active token, it will request one and revoke it immediately afterwards.

# SYNOPSIS

```
somaadm logout [-a|--all]
```

# ARGUMENT TYPES

# PERMISSIONS

This command requires no permissions.

# EXAMPLES

```
./somaadm -u ${USER} logout

./somaadm logout

./somaadm logout -a
```

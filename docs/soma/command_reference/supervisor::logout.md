# DESCRIPTION

This command can be used by a user to revoke active authentication
tokens.

By default it revokes the token currently used by the soma client.
If the all flag is provided, it revokes all tokens for the account that
were issued prior to the revocation.

If the client currently does not have an active token, it will not
request one and revoke it immediately afterwards unless the all flag is
given. In that case it will request a new token and use it to issue the
revocation of all tokens, which includes the new one.

# SYNOPSIS

```
soma logout [-a|--all]
```

# ARGUMENT TYPES

This command takes no arguments.

# PERMISSIONS

This command requires no permissions.

# EXAMPLES

```
soma -u ${USER} logout
soma logout
soma logout -a
soma logout --all
```

# environment definitions

Environments are part of the configuration tree metadata. Within a tree,
every bucket is assigned an environment. Check configurations defined on
repositories can thus be constrained on specific environments.

Environment names must not contain a / character.

# SYNOPSIS OVERVIEW

```
soma environment add ${environment}
soma environment remove ${environment}
soma environment rename ${environment} to ${new-environment}
soma environment list
soma environment show ${environment}
```

See `soma environment help ${command}` for detailed help.

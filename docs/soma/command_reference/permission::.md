# permission definitions

Permissions are a runtime configurable part of permission system, they
can be dynamically created. They allow the execution of the actions that
are mapped to them and can be granted to subjects.

# SYNOPSIS OVERVIEW

```
soma permission add ${permission} to ${category}
soma permission remove [${category}::]${permission} [from ${category}]
soma permission list in ${category}
soma permission show [${category}::]${permission} [in ${category}]
soma permission map ${section}[::${action}] to ${category}::${permission}
soma permission unmap ${section}[::${action}] from ${category}::${permission}
```

See `soma permission help ${command}` for detailed help.

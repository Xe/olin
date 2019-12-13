# Policy Files

Olin offers the ability to customize what the guest webassembly environment can
and cannot access by using a policy file. Here are the options you can put in
one:

```
## This is a comment. All comments require two hashes.

## This is the list of URL matching regexes that the environment is explicitly
## allowed to open. Any resource URL not matching one of these regexes will be
## denied and the webassembly module will be immediately killed.
allow (
  https://mycoolsite.com
  https://discordapp.com/api/webhooks
  random://
  zero://
)

## This is the list of URL matching regexes that the environment is explicitly
## disallowed from opening. Any resource URL matching one of these regexes will
## be denied and the webassembly module will be immediately killed.
disallow (
  https://mycoolsite.com/admin
)

## This is the ram limit in webassembly pages. The default is 128 pages, or
## about eight megabytes.
ram-page-limit 128

## This is the gas limit. This defines how many webassembly instructions the
## environment is allowed to execute. If it goes over this number, the
## webassembly module gets immediately killed. The default is over 33 million
## instructions.
gas-limit 33554432
```

Or more minimally:

```
allow (
  https://mycoolsite.com
  https://discordapp.com/api/webhooks
  random://
  zero://
)

disallow (
  https://mycoolsite.com/admin
)

ram-page-limit 128
gas-limit 33554432
```

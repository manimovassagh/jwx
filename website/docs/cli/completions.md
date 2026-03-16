---
sidebar_position: 4
title: "Shell Completions"
---

# Shell Completions

jwx supports tab completions for Bash, Zsh, Fish, and PowerShell. Completions help you discover commands, flags, and valid values without consulting the documentation.

## Bash

Generate and install Bash completions:

```bash
# System-wide (requires root)
jwx completion bash | sudo tee /etc/bash_completion.d/jwx > /dev/null

# User-only
mkdir -p ~/.local/share/bash-completion/completions
jwx completion bash > ~/.local/share/bash-completion/completions/jwx
```

Restart your shell or source the file to activate:

```bash
source ~/.local/share/bash-completion/completions/jwx
```

## Zsh

Generate and install Zsh completions:

```bash
# Add to your fpath
jwx completion zsh > "${fpath[1]}/_jwx"
```

If `${fpath[1]}` doesn't exist or isn't writable, create a local completions directory:

```bash
mkdir -p ~/.zsh/completions
jwx completion zsh > ~/.zsh/completions/_jwx
```

Then add this to your `~/.zshrc` (before `compinit`):

```bash
fpath=(~/.zsh/completions $fpath)
autoload -Uz compinit && compinit
```

You may need to delete the Zsh completion cache to pick up new completions:

```bash
rm -f ~/.zcompdump
```

## Fish

Generate and install Fish completions:

```bash
jwx completion fish > ~/.config/fish/completions/jwx.fish
```

Fish picks up the completions automatically -- no restart needed.

## PowerShell

Load completions in your current session:

```powershell
jwx completion powershell | Out-String | Invoke-Expression
```

To load completions automatically, add this to your PowerShell profile (`$PROFILE`):

```powershell
jwx completion powershell | Out-String | Invoke-Expression
```

To find your profile path:

```powershell
echo $PROFILE
```

## What gets completed

With completions enabled, pressing Tab will suggest:

- **Subcommands**: `decode`, `sign`, `version`, `completion`
- **Flags**: `--json`, `--clipboard`, `--no-color`, `--alg`, `--secret`, `--key`, `--from`
- **Algorithm values**: When typing after `--alg`, completions suggest valid algorithm names

## Verifying completions work

After installation, test by typing `jwx ` and pressing Tab:

```bash
$ jwx <TAB>
completion  decode  sign  version
```

Or test flag completion:

```bash
$ jwx decode --<TAB>
--clipboard  --json  --no-color
```

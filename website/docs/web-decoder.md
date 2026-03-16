---
sidebar_position: 5
title: Web Decoder
---

# Web Decoder

jwx includes a browser-based JWT decoder at **[manimovassagh.github.io/jwx](https://manimovassagh.github.io/jwx/)**. It provides the same decoding functionality as the CLI, with a visual split-pane interface.

## Features

- **Instant decoding** -- tokens are decoded as you type, with no button to click
- **Syntax highlighting** -- header, payload, and signature are color-coded for readability
- **Algorithm badge** -- shows the signing algorithm (HS256, RS256, etc.) at a glance
- **Expiry detection** -- displays whether the token is valid or expired with a colored badge
- **Timestamp conversion** -- `iat`, `exp`, `nbf`, and `auth_time` claims are shown with human-readable dates and relative times
- **Copy sections** -- copy the header, payload, or signature JSON with one click
- **Sample tokens** -- pre-loaded sample tokens to explore different JWT structures
- **Dark and light themes** -- toggle between themes, or let it follow your OS preference
- **Shareable URLs** -- the token is stored in the URL hash, so you can share decoded views with teammates
- **Mobile responsive** -- works on phones and tablets with a stacked layout

## Privacy model

The web decoder runs entirely in your browser:

- **No server-side processing** -- all decoding happens in JavaScript on your device
- **No network requests** -- tokens are never sent anywhere (enforced by Content Security Policy)
- **No tracking or analytics** -- no cookies, no local storage (except theme preference), no third-party scripts
- **No accounts** -- completely anonymous, no sign-up required

The page includes a strict Content Security Policy that prevents any external network connections for scripts or data.

## How to use it

1. Open [manimovassagh.github.io/jwx](https://manimovassagh.github.io/jwx/)
2. Paste your JWT token in the left panel
3. The right panel instantly shows the decoded header, payload, and signature
4. Click a sample token pill at the bottom of the left panel to see example tokens

## Sharing decoded tokens

When you paste a token, the URL hash updates to include the encoded token. You can copy the URL and share it with teammates -- they'll see the same decoded view when they open it.

Example: `https://manimovassagh.github.io/jwx/#eyJhbGciOiJIUzI1NiIs...`

:::caution
Only share tokens that are safe to share. Even though the token stays in the URL fragment (which is not sent to the server), anyone with the link can decode the token's payload. Never share production tokens containing sensitive data.
:::

## Keyboard accessibility

The web decoder is fully keyboard-accessible:

- **Tab** navigates between the input area, sample pills, copy buttons, and theme toggle
- **Skip link** at the top of the page jumps directly to the input area
- **Focus indicators** are visible on all interactive elements

## When to use the CLI vs. the web decoder

| Use the CLI when... | Use the web decoder when... |
|---|---|
| You're already in the terminal | You want a visual, split-pane view |
| You need to pipe output to other tools | You want to share a decoded token via URL |
| You need JSON output for scripts | You want to quickly try sample tokens |
| You need to sign tokens | You're on a machine where you can't install software |
| You're working with sensitive production tokens | You're debugging with test/dev tokens |

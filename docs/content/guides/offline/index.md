---
title: "Using CannaNote Offline"
description: "How offline mode works, what to expect, and how to keep your data safe."
sidebar_label: "Offline"
sidebar_order: 1
section: "guides"
keywords: ["offline", "data", "storage", "cache", "privacy", "local-first"]
related_pages: ["getting-started/index"]
---

# Using CannaNote Offline

CannaNote is built offline-first. Your journal entries live on your device, not our servers.

---

## How It Works

### Your Data Lives on Your Device

When you log a session, it's stored directly in your browser's local database. This means:

- **No internet required** to log sessions or view history
- **Instant access** - no server round-trips
- **Complete privacy** - consumption data stays on your device
- **Works anywhere** - log sessions with no cell service

### Pages Save When You Visit

Documentation, guides, and educational content use a "cache on visit" approach:

1. **First visit** - page downloads and saves automatically
2. **Future visits** - loads instantly from your device
3. **Updates** - refreshes when you're online

This keeps the app fast while giving you offline access to what matters.

---

## What Works Offline

### Always Available

- ✅ Your journal dashboard
- ✅ Logging new sessions
- ✅ Viewing session history
- ✅ Searching your entries
- ✅ Pages you've visited before

### Requires Internet

- ❌ First-time account creation
- ❌ Signing in on a new device
- ❌ Cloud sync (if enabled)
- ❌ Pages you haven't visited yet

---

## Preparing for Offline

### Save Pages for Later

1. Connect to the internet
2. Visit the pages you want available offline
3. Done - they're cached automatically

No download button needed. Your browser remembers where you've been.

### Install as an App

For the best offline experience, install CannaNote on your device:

**Mobile:** Visit cannanote.app → "Add to Home Screen"

**Desktop:** Visit cannanote.app → Click install icon in address bar

Installing ensures core features stay cached and ready.

---

## Keeping Your Data Safe

### Cache vs. Site Data

This distinction matters when troubleshooting:

| Storage Type | Contains | "Clear Cache" Removes? | "Clear Site Data" Removes? |
|--------------|----------|------------------------|----------------------------|
| **Cache** | Pages, CSS, images | Yes | Yes |
| **Site Data** | Your journal entries | No | **Yes** |

### The Warning

When browsers offer "Clear browsing data," be careful:

- **"Clear cache"** = Safe. Just re-downloads pages.
- **"Clear site data"** or **"Clear all data"** = **Deletes your journal.**

> Before clearing site data, export your sessions first. Better safe than sorry.

---

## Browser Guides

Each browser handles offline storage a bit differently. Find yours:

- [Chrome](chrome) - Full PWA support, reliable offline
- [Safari](safari) - Full support, install via Share menu
- [Firefox](firefox) - Full support, strong privacy defaults
- [Brave](brave) - Full support, built-in privacy features
- [DuckDuckGo](duckduckgo) - Full support, watch the Fire Button
- [Zen](zen) - Full support, Firefox-based

### Coming Soon

**Ladybird** - We're watching this independent browser project closely. Once it reaches stable release, we'll add a guide.

---

## Common Issues

**"Page not available offline"**
→ Visit that page while online. It'll cache automatically.

**"Showing outdated content"**
→ Connect to internet and refresh. Cache updates automatically.

**"Can't install as app"**
→ Try Chrome, Brave, or Safari. Some browsers don't support PWA installation.

**"Lost data after clearing browser"**
→ If you had cloud sync, sign back in. If not, the data may be unrecoverable.

For detailed troubleshooting, see your [browser's guide](#browser-guides).

---

## Need Help?

- **Questions?** [support@cannanote.app](mailto:support@cannanote.app)
- **Found a bug?** [Report on GitHub](https://github.com/kcanamar/cannanote/issues)

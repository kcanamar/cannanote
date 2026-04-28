---
title: "Safari Offline Guide"
description: "How to use CannaNote offline in Safari, manage cached data, and add to your Home Screen."
sidebar_label: "Safari"
sidebar_order: 3
section: "guides"
keywords: ["safari", "offline", "pwa", "cache", "apple", "ios", "macos"]
related_pages: ["guides/offline/index"]
---

# Safari Offline Guide

Safari on macOS and iOS fully supports CannaNote's offline features. Here's what you need to know.

---

## Installing CannaNote

### On iPhone / iPad

1. Visit **cannanote.app** in Safari
2. Tap the **Share** button (square with arrow)
3. Scroll down and tap **Add to Home Screen**
4. Tap **Add**
5. CannaNote now appears as an app icon

### On Mac

1. Visit **cannanote.app** in Safari
2. Click **File** in the menu bar
3. Click **Add to Dock** (macOS Sonoma+) or use the Share menu
4. CannaNote opens in its own window

---

## Managing Cached Data

### Clear Cache Only (Safe)

This removes downloaded pages but **keeps your journal entries**.

**On Mac:**
1. Enable the Develop menu (see below if needed)
2. Click **Develop** in the menu bar
3. Click **Empty Caches**

**On iPhone / iPad:**
1. Go to **Settings** → **Safari**
2. Tap **Clear History and Website Data**
3. ⚠️ **Warning**: This clears everything. See "Safari Limitations" below.

### Enable the Develop Menu (Mac)

You'll need this for cache management and debugging.

1. Open **Safari** → **Settings** (or Preferences)
2. Click the **Advanced** tab
3. Check **"Show Develop menu in menu bar"**

### View Stored Data

**On Mac:**
1. Go to **Safari** → **Settings**
2. Click the **Privacy** tab
3. Click **Manage Website Data**
4. Search for **cannanote.app**

**On iPhone / iPad:**
1. Go to **Settings** → **Safari**
2. Scroll to **Advanced**
3. Tap **Website Data**
4. Search for **cannanote.app**

---

## Developer Tools

### Open Web Inspector (Mac)

1. Enable the Develop menu (see above)
2. Press **Cmd + Option + I**, or
3. Click **Develop** → **Show Web Inspector**

### View Service Worker

1. Open Web Inspector
2. Click the **Storage** tab (or Application)
3. Look for **Service Workers** in the sidebar

### On iPhone / iPad

Safari's developer tools require connecting to a Mac:

1. On your iPhone/iPad: **Settings** → **Safari** → **Advanced** → Enable **Web Inspector**
2. Connect your device to a Mac via USB
3. On Mac: Open Safari → **Develop** menu → Select your device

---

## Safari-Specific Notes

### Limitations

- **iOS cache clearing is all-or-nothing**: Unlike desktop browsers, iOS Safari doesn't let you clear just the cache. "Clear History and Website Data" removes everything.
- **Storage limits**: Safari may clear old cached data if your device runs low on storage. Install as an app for more reliable persistence.
- **7-day limit (older iOS)**: iOS versions before 16 could delete service worker data after 7 days of no use. Keep using the app regularly, or update to iOS 16+.

### Tips for iPhone/iPad

- **Add to Home Screen** for the most reliable offline experience
- **Use regularly** to prevent iOS from cleaning up storage
- **Enable cloud sync** as a backup if you're concerned about data loss

---

## Quick Reference

| Action | How |
|--------|-----|
| Install (iOS) | Share → Add to Home Screen |
| Install (Mac) | File → Add to Dock |
| Clear cache (Mac) | Develop → Empty Caches |
| View stored data | Settings → Privacy → Manage Website Data |
| Open Web Inspector | Cmd + Option + I (with Develop menu enabled) |

---

[← Back to Offline Guide](index)

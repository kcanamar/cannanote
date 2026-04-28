---
title: "Brave Offline Guide"
description: "How to use CannaNote offline in Brave, manage cached data, and install as a PWA."
sidebar_label: "Brave"
sidebar_order: 5
section: "guides"
keywords: ["brave", "offline", "pwa", "cache", "privacy"]
related_pages: ["guides/offline/index"]
---

# Brave Offline Guide

Brave is Chromium-based with built-in privacy features. CannaNote works great here.

---

## Installing CannaNote

1. Visit **cannanote.app** in Brave
2. Look for the **install icon** (⊕) in the address bar
3. Click **Install**
4. CannaNote opens in its own window

Same process as Chrome - Brave has full PWA support.

---

## Managing Cached Data

### Clear Cache Only (Safe)

This removes downloaded pages but **keeps your journal entries**.

1. Click the **☰** menu (or **⋮** depending on version)
2. Go to **Settings**
3. Click **Privacy and security**
4. Click **Clear browsing data**
5. Select **"Cached images and files"** only
6. Make sure "Cookies and other site data" is **unchecked**
7. Click **Clear data**

### View Stored Data

1. Click the menu icon
2. Go to **Settings**
3. Click **Privacy and security**
4. Click **Site and Shields Settings**
5. Click **View permissions and data stored**
6. Search for **cannanote.app**

---

## Brave Shields

Brave's Shields block trackers and ads by default. Good news: this doesn't affect CannaNote.

### Why It Works Fine

- CannaNote uses no third-party trackers
- All storage is first-party (your data stays on your device)
- No ads, no analytics scripts to block

### If Something Seems Off

You can check Shield settings for any site:

1. Click the **Brave Shield** icon (lion) in the address bar
2. Verify shields are working normally
3. CannaNote should work with shields up or down

---

## Developer Tools

Brave uses Chromium DevTools, same as Chrome.

### Open DevTools

Press **F12** or **Ctrl+Shift+I** (Cmd+Option+I on Mac)

### View Service Worker

1. Open DevTools (F12)
2. Click the **Application** tab
3. In the sidebar, click **Service Workers**

### Unregister Service Worker

1. Open DevTools (F12)
2. Go to **Application** → **Service Workers**
3. Click **Unregister**
4. Refresh the page

---

## Brave-Specific Notes

### Brave Rewards

CannaNote doesn't participate in Brave Rewards. We don't show ads, so there's nothing to tip or earn from.

### Private Window with Tor

Brave offers "Private Window with Tor." While this works, it's overkill for CannaNote and adds latency. Regular browsing or a standard private window works fine.

### Sync

Brave Sync handles bookmarks and settings but not website storage. For cross-device journal access, use CannaNote's cloud sync feature.

---

## Quick Reference

| Action | How |
|--------|-----|
| Install as app | Address bar → ⊕ icon → Install |
| Clear cache safely | Settings → Privacy → Clear browsing data → Select "Cached images and files" only |
| View stored data | Settings → Privacy → Site and Shields Settings → View permissions and data |
| Check Shields | Click lion icon in address bar |
| Open DevTools | F12 or Ctrl+Shift+I |
| Unregister SW | DevTools → Application → Service Workers → Unregister |

---

[← Back to Offline Guide](index)

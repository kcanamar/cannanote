---
title: "Firefox Offline Guide"
description: "How to use CannaNote offline in Firefox, manage cached data, and install as an app."
sidebar_label: "Firefox"
sidebar_order: 4
section: "guides"
keywords: ["firefox", "offline", "pwa", "cache", "mozilla"]
related_pages: ["guides/offline/index"]
---

# Firefox Offline Guide

Firefox offers strong privacy defaults and full support for CannaNote's offline features.

---

## Installing CannaNote

### On Desktop

Firefox supports installing web apps:

1. Visit **cannanote.app** in Firefox
2. Look for the **install icon** in the address bar
3. Click to install
4. CannaNote opens in its own window

*Note: If you don't see an install option, Firefox may require an extension for full PWA support. The web version works great either way.*

### On Android

1. Visit **cannanote.app** in Firefox
2. Tap the **⋮** menu
3. Tap **Add to Home screen**
4. CannaNote appears as an app icon

---

## Managing Cached Data

### Clear Cache Only (Safe)

This removes downloaded pages but **keeps your journal entries**.

1. Click the **☰** menu (top right)
2. Click **Settings**
3. Go to **Privacy & Security**
4. Scroll to **Cookies and Site Data**
5. Click **Clear Data**
6. Check **"Cached Web Content"** only
7. Make sure "Cookies and Site Data" is **unchecked**
8. Click **Clear**

### View Stored Data

See what CannaNote stores on your device.

1. Click the **☰** menu
2. Click **Settings**
3. Go to **Privacy & Security**
4. Scroll to **Cookies and Site Data**
5. Click **Manage Data**
6. Search for **cannanote.app**

---

## Developer Tools

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

### View Cached Files

1. Open DevTools (F12)
2. Go to **Application** → **Cache Storage**
3. Expand to see cached pages and assets

---

## Firefox-Specific Notes

### Enhanced Tracking Protection

Firefox's privacy features don't interfere with CannaNote because:
- We don't use third-party trackers
- All data storage is first-party
- No cross-site cookies

### Private Browsing

CannaNote works in private windows, but:
- Data is deleted when you close the window
- Not recommended for regular journaling
- Use normal browsing for persistent storage

### Firefox Containers

If you use Multi-Account Containers:
- Each container has separate storage
- Your journal data stays in the container where you created it
- Stick to one container for CannaNote

---

## Quick Reference

| Action | How |
|--------|-----|
| Install as app | Address bar install icon (if available) |
| Add to Home (Android) | Menu → Add to Home screen |
| Clear cache safely | Settings → Privacy → Clear Data → Select "Cached Web Content" only |
| View stored data | Settings → Privacy → Manage Data |
| Open DevTools | F12 or Ctrl+Shift+I |
| Unregister SW | DevTools → Application → Service Workers → Unregister |

---

[← Back to Offline Guide](index)

---
title: "Chrome Offline Guide"
description: "How to use CannaNote offline in Google Chrome, manage cached data, and install as a PWA."
sidebar_label: "Chrome"
sidebar_order: 2
section: "guides"
keywords: ["chrome", "offline", "pwa", "cache", "google"]
related_pages: ["guides/offline/index"]
---

# Chrome Offline Guide

Chrome offers excellent PWA support and reliable offline functionality. Here's everything you need to know.

---

## Installing CannaNote

Installing as an app gives you the best offline experience.

1. Visit **cannanote.app** in Chrome
2. Look for the **install icon** (⊕) in the address bar
3. Click **Install**
4. CannaNote now opens in its own window

Once installed, core features stay cached automatically.

---

## Managing Cached Data

### Clear Cache Only (Safe)

This removes downloaded pages but **keeps your journal entries**.

1. Click the **⋮** menu (top right)
2. Go to **Settings**
3. Click **Privacy and security**
4. Click **Clear browsing data**
5. Select **"Cached images and files"** only
6. Make sure "Cookies and other site data" is **unchecked**
7. Click **Clear data**

### View Stored Data

See exactly what CannaNote stores on your device.

1. Click the **⋮** menu
2. Go to **Settings**
3. Click **Privacy and security**
4. Click **Site settings**
5. Click **View permissions and data stored**
6. Search for **cannanote.app**

---

## Developer Tools

For troubleshooting or if support asks you to check something:

### View Service Worker

1. Press **F12** to open DevTools
2. Click the **Application** tab
3. In the sidebar, click **Service Workers**
4. You'll see CannaNote's service worker status

### Unregister Service Worker

If something's broken and you want a fresh start:

1. Open DevTools (F12)
2. Go to **Application** → **Service Workers**
3. Click **Unregister** next to cannanote.app
4. Refresh the page

### View Cached Files

1. Open DevTools (F12)
2. Go to **Application** → **Cache Storage**
3. Expand to see all cached pages and assets

---

## Chrome-Specific Notes

- **Sync across devices**: Chrome can sync your browsing data, but CannaNote's local storage is device-specific. Use our cloud sync for cross-device access.
- **Incognito mode**: CannaNote works in incognito, but data is cleared when you close the window. Not recommended for regular use.
- **Mobile Chrome (Android)**: Same features as desktop. Install via "Add to Home screen" in the menu.

---

## Quick Reference

| Action | How |
|--------|-----|
| Install as app | Address bar → ⊕ icon → Install |
| Clear cache safely | Settings → Privacy → Clear browsing data → Select "Cached images and files" only |
| View stored data | Settings → Privacy → Site settings → View permissions and data |
| Open DevTools | F12 or Ctrl+Shift+I (Cmd+Option+I on Mac) |
| Unregister SW | DevTools → Application → Service Workers → Unregister |

---

[← Back to Offline Guide](index)

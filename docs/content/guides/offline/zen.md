---
title: "Zen Browser Offline Guide"
description: "How to use CannaNote offline in Zen Browser, manage cached data, and install as an app."
sidebar_label: "Zen"
sidebar_order: 7
section: "guides"
keywords: ["zen", "offline", "pwa", "cache", "firefox", "privacy"]
related_pages: ["guides/offline/index"]
---

# Zen Browser Offline Guide

Zen is a privacy-focused browser built on Firefox. Since it uses the Gecko engine, most Firefox instructions apply here too.

---

## Installing CannaNote

1. Visit **cannanote.app** in Zen
2. Look for an **install icon** in the address bar or menu
3. Follow prompts to install
4. CannaNote opens in its own window

*Note: Zen's PWA support follows Firefox's implementation. If install options aren't visible, the web version works great.*

---

## Managing Cached Data

### Clear Cache Only (Safe)

This removes downloaded pages but **keeps your journal entries**.

1. Open **Settings** (usually via ☰ menu)
2. Go to **Privacy & Security**
3. Find **Cookies and Site Data**
4. Click **Clear Data**
5. Check **"Cached Web Content"** only
6. Make sure "Cookies and Site Data" is **unchecked**
7. Click **Clear**

### View Stored Data

1. Open **Settings**
2. Go to **Privacy & Security**
3. Find **Cookies and Site Data**
4. Click **Manage Data**
5. Search for **cannanote.app**

---

## Developer Tools

Zen uses Firefox's developer tools.

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

## Zen-Specific Notes

### Privacy Features

Zen emphasizes privacy and customization. Like Firefox, its privacy features don't interfere with CannaNote:

- We use no third-party trackers
- All storage is first-party
- Your data stays on your device

### Customization

Zen is known for its customizable interface. CannaNote adapts to your system theme (light/dark), so it should blend with however you've configured Zen.

### Updates

Since Zen is Firefox-based, it benefits from Firefox's security updates. Keep your browser updated for the best experience.

### Community

Zen has an active community. If you run into browser-specific issues, their community channels may have insights. For CannaNote-specific issues, reach out to us directly.

---

## Quick Reference

| Action | How |
|--------|-----|
| Install as app | Address bar install icon (if available) |
| Clear cache safely | Settings → Privacy → Clear Data → Select "Cached Web Content" only |
| View stored data | Settings → Privacy → Manage Data |
| Open DevTools | F12 or Ctrl+Shift+I |
| Unregister SW | DevTools → Application → Service Workers → Unregister |

---

[← Back to Offline Guide](index)

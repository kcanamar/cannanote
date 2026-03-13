// CannaNote Session Notifications
// Handles Web Notification scheduling and background notification checks
// Zero dependencies, integrates with session-db.js

const SessionNotifications = (function() {
  // Check interval for pending notifications (when page is open)
  const CHECK_INTERVAL_MS = 30 * 1000; // 30 seconds

  let checkIntervalId = null;
  let isInitialized = false;

  // Request notification permission
  async function requestPermission() {
    if (!('Notification' in window)) {
      console.log('[Notifications] Not supported in this browser');
      return false;
    }

    if (Notification.permission === 'granted') {
      return true;
    }

    if (Notification.permission === 'denied') {
      console.log('[Notifications] Permission denied');
      return false;
    }

    const permission = await Notification.requestPermission();
    return permission === 'granted';
  }

  // Initialize notification system
  async function init() {
    if (isInitialized) return;

    // Start checking for pending notifications
    startChecking();

    // Listen for visibility changes to pause/resume checking
    document.addEventListener('visibilitychange', () => {
      if (document.visibilityState === 'visible') {
        startChecking();
        checkPendingNotifications(); // Immediate check when tab becomes visible
      } else {
        stopChecking();
      }
    });

    isInitialized = true;
    console.log('[Notifications] Initialized');
  }

  // Start periodic checking
  function startChecking() {
    if (checkIntervalId) return;
    checkIntervalId = setInterval(checkPendingNotifications, CHECK_INTERVAL_MS);
    console.log('[Notifications] Started checking');
  }

  // Stop periodic checking
  function stopChecking() {
    if (checkIntervalId) {
      clearInterval(checkIntervalId);
      checkIntervalId = null;
      console.log('[Notifications] Stopped checking');
    }
  }

  // Check for and show pending notifications
  async function checkPendingNotifications() {
    if (!window.SessionDB) {
      console.log('[Notifications] SessionDB not available');
      return;
    }

    try {
      const pending = await SessionDB.getPendingNotifications();

      for (const notification of pending) {
        // Show the notification
        await showNotification(notification);

        // Mark as shown
        await SessionDB.markNotificationShown(notification.id);
      }
    } catch (err) {
      console.error('[Notifications] Check failed:', err);
    }
  }

  // Show a notification
  async function showNotification(notificationData) {
    const hasPermission = await requestPermission();
    if (!hasPermission) return;

    try {
      // Get session details for the notification
      const session = await SessionDB.get('sessions', notificationData.sessionId);
      if (!session || session.status !== 'active') {
        // Session no longer active, don't show notification
        return;
      }

      let product = null;
      if (session.productId) {
        product = await SessionDB.get('products', session.productId);
      }

      const strainName = product?.strainName || 'Your session';
      const elapsed = Math.floor((Date.now() - session.timestamp) / 60000);
      let elapsedText = `${elapsed} minutes`;
      if (elapsed >= 60) {
        const hours = Math.floor(elapsed / 60);
        const mins = elapsed % 60;
        elapsedText = `${hours}h ${mins}m`;
      }

      const notification = new Notification('CannaNote Check-in', {
        body: `${strainName} - ${elapsedText} since your session. How are you feeling?`,
        icon: '/assets/images/favicon/favicon-96x96.png',
        tag: `checkin-${session.id}`,
        requireInteraction: true,
        data: {
          sessionId: session.id,
          url: `/app/sessions/${session.id}/checkin`
        }
      });

      notification.onclick = function(event) {
        event.preventDefault();
        window.focus();
        window.location.href = this.data.url;
        this.close();
      };

      console.log('[Notifications] Shown for session:', session.id);

    } catch (err) {
      console.error('[Notifications] Failed to show:', err);
    }
  }

  // Schedule a notification for a session
  async function scheduleCheckIn(sessionId, intervalMinutes) {
    const scheduledTime = Date.now() + (intervalMinutes * 60 * 1000);

    try {
      await SessionDB.scheduleNotification(sessionId, scheduledTime);
      console.log('[Notifications] Scheduled for', intervalMinutes, 'minutes');

      // Also set a local timeout as backup (for when page stays open)
      setTimeout(async () => {
        const session = await SessionDB.get('sessions', sessionId);
        if (session && session.status === 'active') {
          await showNotification({ sessionId, scheduledTime });
        }
      }, intervalMinutes * 60 * 1000);

      return true;
    } catch (err) {
      console.error('[Notifications] Schedule failed:', err);
      return false;
    }
  }

  // Cancel all notifications for a session
  async function cancelSessionNotifications(sessionId) {
    try {
      await SessionDB.clearSessionNotifications(sessionId);
      console.log('[Notifications] Cancelled for session:', sessionId);
    } catch (err) {
      console.error('[Notifications] Cancel failed:', err);
    }
  }

  // Get notification permission status
  function getPermissionStatus() {
    if (!('Notification' in window)) {
      return 'unsupported';
    }
    return Notification.permission;
  }

  // Check if notifications are supported
  function isSupported() {
    return 'Notification' in window;
  }

  // Public API
  return {
    init,
    requestPermission,
    scheduleCheckIn,
    cancelSessionNotifications,
    checkPendingNotifications,
    getPermissionStatus,
    isSupported
  };
})();

// Initialize on load
if (typeof window !== 'undefined') {
  window.SessionNotifications = SessionNotifications;

  document.addEventListener('DOMContentLoaded', () => {
    // Only init if SessionDB is available
    if (window.SessionDB) {
      SessionNotifications.init();
    } else {
      // Wait for SessionDB to be ready
      const checkForDB = setInterval(() => {
        if (window.SessionDB) {
          clearInterval(checkForDB);
          SessionNotifications.init();
        }
      }, 100);

      // Give up after 5 seconds
      setTimeout(() => clearInterval(checkForDB), 5000);
    }
  });
}

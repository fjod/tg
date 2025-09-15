/**
 * Telegram Web App utilities for handling SDK integration and authentication
 */

class TelegramWebApp {
  constructor() {
    this.tg = window.Telegram?.WebApp;
    this.user = null;
    this.initData = null;
  }

  /**
   * Initialize Telegram Web App
   */
  init() {
    if (!this.tg) {
      console.warn('Telegram WebApp SDK not found');
      return false;
    }

    // Expand the app to full height
    this.tg.expand();

    // Enable closing confirmation
    this.tg.enableClosingConfirmation();

    // Set header color to match theme
    this.tg.setHeaderColor('bg_color');

    // Get user data and init data for authentication
    this.user = this.tg.initDataUnsafe?.user;
    this.initData = this.tg.initData;

    console.log('Telegram WebApp initialized', {
      user: this.user,
      platform: this.tg.platform,
      version: this.tg.version,
      colorScheme: this.tg.colorScheme,
    });

    return true;
  }

  /**
   * Get the authentication data for API calls
   * @returns {string|null} The initData string for authentication
   */
  getAuthData() {
    return this.initData;
  }

  /**
   * Get current user information
   * @returns {Object|null} User object with id, first_name, username, etc.
   */
  getUser() {
    return this.user;
  }

  /**
   * Check if app is running inside Telegram
   * @returns {boolean}
   */
  isInTelegram() {
    return !!this.tg && !!this.initData;
  }

  /**
   * Detect the platform where the Mini App is running
   * @returns {string} Platform type: 'web', 'desktop', 'mobile', or 'unknown'
   */
  detectPlatform() {
    if (!this.tg) {
      return 'unknown';
    }

    const platform = this.tg.platform;

    // Platform detection based on Telegram WebApp platform property
    if (platform) {
      // Common platform values: 'web', 'tdesktop', 'ios', 'android', 'macos'
      if (platform === 'web') {
        return 'web';
      } else if (platform === 'tdesktop' || platform === 'macos') {
        return 'desktop';
      } else if (platform === 'ios' || platform === 'android') {
        return 'mobile';
      }
    }

    // Fallback detection based on user agent
    const userAgent = navigator.userAgent.toLowerCase();
    if (userAgent.includes('electron')) {
      return 'desktop';
    } else if (userAgent.includes('mobile') || userAgent.includes('android') || userAgent.includes('iphone')) {
      return 'mobile';
    }

    // Default to web if unclear
    return 'web';
  }

  /**
   * Check if running in native Telegram client (desktop or mobile)
   * @returns {boolean}
   */
  isNativeClient() {
    const platform = this.detectPlatform();
    return platform === 'desktop' || platform === 'mobile';
  }

  /**
   * Show main button with callback
   * @param {string} text - Button text
   * @param {Function} callback - Click callback
   */
  showMainButton(text, callback) {
    if (!this.tg) return;

    this.tg.MainButton.setText(text);
    this.tg.MainButton.onClick(callback);
    this.tg.MainButton.show();
  }

  /**
   * Hide main button
   */
  hideMainButton() {
    if (!this.tg) return;
    this.tg.MainButton.hide();
  }

  /**
   * Show popup alert
   * @param {string} message - Alert message
   * @param {Function} callback - Callback after OK
   */
  showAlert(message, callback) {
    if (!this.tg) {
      alert(message);
      if (callback) callback();
      return;
    }

    this.tg.showAlert(message, callback);
  }

  /**
   * Show confirm dialog
   * @param {string} message - Confirm message
   * @param {Function} callback - Callback with boolean result
   */
  showConfirm(message, callback) {
    if (!this.tg) {
      const result = confirm(message);
      if (callback) callback(result);
      return;
    }

    this.tg.showConfirm(message, callback);
  }

  /**
   * Get theme colors
   * @returns {Object} Theme parameters
   */
  getTheme() {
    if (!this.tg) {
      return {
        bg_color: '#ffffff',
        text_color: '#000000',
        hint_color: '#999999',
        link_color: '#007aff',
        button_color: '#007aff',
        button_text_color: '#ffffff',
      };
    }

    return this.tg.themeParams;
  }

  /**
   * Close the mini app
   */
  close() {
    if (this.tg) {
      this.tg.close();
    }
  }

  /**
   * Send haptic feedback
   * @param {string} type - 'impact', 'notification', or 'selection'
   * @param {string} style - For impact: 'light', 'medium', 'heavy'; For notification: 'error', 'success', 'warning'
   */
  hapticFeedback(type, style = 'light') {
    if (!this.tg?.HapticFeedback) return;

    switch (type) {
      case 'impact':
        this.tg.HapticFeedback.impactOccurred(style);
        break;
      case 'notification':
        this.tg.HapticFeedback.notificationOccurred(style);
        break;
      case 'selection':
        this.tg.HapticFeedback.selectionChanged();
        break;
    }
  }

  /**
   * Open Telegram link (redirect to message or other Telegram content)
   * @param {string} url - Telegram URL (e.g., t.me/c/chatid/messageid)
   * @param {Object} options - Additional options
   */
  openTelegramLink(url, options = {}) {
    if (!url) {
      return false;
    }

    if (this.tg?.openTelegramLink) {
      try {
        this.tg.openTelegramLink(url, options);
        return true;
      } catch (error) {
        return false;
      }
    } else {
      // Fallback for development or if WebApp API is unavailable
      try {
        window.open(url, '_blank');
        return true;
      } catch (error) {
        return false;
      }
    }
  }

  /**
   * Generate Telegram message URL for redirection
   * @param {number} messageId - Telegram message ID
   * @param {number} chatId - Chat ID (optional, uses current user if not provided)
   * @returns {string} Telegram message URL
   */
  generateMessageUrl(messageId, chatId = null) {
    if (!messageId) {
      return null;
    }

    // For web clients, use username-based URL format
    if (this.user?.username) {
      return `https://t.me/${this.user.username}/${messageId}`;
    }
    
    // Fallback: Try the original format
    const targetChatId = chatId || this.user?.id;
    
    if (!targetChatId) {
      return null;
    }

    return `https://t.me/c/${targetChatId}/${messageId}`;
  }

  /**
   * Redirect to a specific Telegram message
   * @param {Object} message - Message object with telegram_message_id
   * @param {number} chatId - Optional chat ID, uses current user if not provided
   * @returns {boolean} True if redirection was attempted
   */
  redirectToMessage(message, chatId = null) {
    if (!message || !message.telegram_message_id) {
      this.showAlert('Unable to redirect to message: Invalid message data');
      return false;
    }

    const platform = this.detectPlatform();
    
    // Only attempt redirect in web view, show helpful message for other platforms
    if (platform !== 'web') {
      this.showAlert('💡 Message redirection to Saved Messages only works in Telegram Web.\n\nTo view this message:\n1. Open Telegram in your web browser\n2. Go to Saved Messages\n3. Look for message ID: ' + message.telegram_message_id);
      return false;
    }

    // For web clients, attempt redirection
    try {
      const messageUrl = this.generateMessageUrl(message.telegram_message_id, chatId);
      
      if (!messageUrl) {
        this.showAlert('Unable to generate message link');
        return false;
      }

      const success = this.openTelegramLink(messageUrl);
      
      if (success) {
        this.hapticFeedback('impact', 'medium');
      } else {
        this.showAlert('Unable to open message. Please navigate to Saved Messages manually.');
      }

      return success;
    } catch (error) {
      this.showAlert('Unable to open message in Telegram. Please try again.');
      this.hapticFeedback('notification', 'error');
      return false;
    }
  }

}

// Create singleton instance
const telegramApp = new TelegramWebApp();

export default telegramApp;
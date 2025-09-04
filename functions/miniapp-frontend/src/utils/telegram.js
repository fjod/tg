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
}

// Create singleton instance
const telegramApp = new TelegramWebApp();

export default telegramApp;
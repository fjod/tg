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

  /**
   * Open Telegram link (redirect to message or other Telegram content)
   * @param {string} url - Telegram URL (e.g., t.me/c/chatid/messageid)
   * @param {Object} options - Additional options
   */
  openTelegramLink(url, options = {}) {
    console.log('🔵 TelegramApp: openTelegramLink called');
    console.log('🔵 TelegramApp: URL:', url);
    console.log('🔵 TelegramApp: Options:', options);
    console.log('🔵 TelegramApp: this.tg exists:', !!this.tg);
    console.log('🔵 TelegramApp: this.tg.openTelegramLink exists:', !!this.tg?.openTelegramLink);
    
    if (!url) {
      console.error('🔴 TelegramApp: No URL provided for Telegram link');
      return false;
    }

    if (this.tg?.openTelegramLink) {
      console.log('🔵 TelegramApp: Using Telegram WebApp openTelegramLink method');
      try {
        this.tg.openTelegramLink(url, options);
        console.log('🟢 TelegramApp: openTelegramLink called successfully');
        return true;
      } catch (error) {
        console.error('🔴 TelegramApp: Failed to open Telegram link via WebApp:', error);
        return false;
      }
    } else {
      console.warn('🟡 TelegramApp: Telegram WebApp openTelegramLink not available, using fallback');
      console.log('🔵 TelegramApp: WebApp object:', this.tg);
      
      // Fallback for development or if WebApp API is unavailable
      try {
        console.log('🔵 TelegramApp: Using window.open fallback');
        window.open(url, '_blank');
        console.log('🟢 TelegramApp: window.open called successfully');
        return true;
      } catch (error) {
        console.error('🔴 TelegramApp: Failed to open Telegram link via fallback:', error);
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
    console.log('🔵 TelegramApp: generateMessageUrl called');
    console.log('🔵 TelegramApp: Input messageId:', messageId);
    console.log('🔵 TelegramApp: Input chatId:', chatId);
    console.log('🔵 TelegramApp: this.user:', this.user);
    console.log('🔵 TelegramApp: this.user?.id:', this.user?.id);
    
    if (!messageId) {
      console.error('🔴 TelegramApp: Message ID is required for URL generation');
      return null;
    }

    // Use provided chatId or fallback to current user
    const targetChatId = chatId || this.user?.id;
    console.log('🔵 TelegramApp: targetChatId resolved to:', targetChatId);
    
    if (!targetChatId) {
      console.error('🔴 TelegramApp: No chat ID available for URL generation');
      console.error('🔴 TelegramApp: chatId provided:', chatId);
      console.error('🔴 TelegramApp: this.user?.id:', this.user?.id);
      return null;
    }

    // For private chats, use the direct message format
    // Note: This format may need adjustment based on actual chat structure
    const url = `https://t.me/c/${targetChatId}/${messageId}`;
    console.log('🔵 TelegramApp: Generated message URL:', url);
    
    return url;
  }

  /**
   * Redirect to a specific Telegram message
   * @param {Object} message - Message object with telegram_message_id
   * @param {number} chatId - Optional chat ID, uses current user if not provided
   * @returns {boolean} True if redirection was attempted
   */
  redirectToMessage(message, chatId = null) {
    console.log('🔵 TelegramApp: redirectToMessage called');
    console.log('🔵 TelegramApp: Input message:', message);
    console.log('🔵 TelegramApp: Input chatId:', chatId);
    console.log('🔵 TelegramApp: Current user:', this.user);
    console.log('🔵 TelegramApp: Telegram WebApp available:', !!this.tg);
    console.log('🔵 TelegramApp: InitData available:', !!this.initData);
    
    if (!message || !message.telegram_message_id) {
      console.error('🔴 TelegramApp: Invalid message object for redirection');
      console.error('🔴 TelegramApp: Message exists:', !!message);
      console.error('🔴 TelegramApp: telegram_message_id exists:', !!message?.telegram_message_id);
      this.showAlert('Unable to redirect to message: Invalid message data');
      return false;
    }

    try {
      console.log('🔵 TelegramApp: Generating message URL...');
      const messageUrl = this.generateMessageUrl(message.telegram_message_id, chatId);
      console.log('🔵 TelegramApp: Generated URL:', messageUrl);
      
      if (!messageUrl) {
        console.error('🔴 TelegramApp: URL generation failed');
        throw new Error('Failed to generate message URL');
      }

      console.log('🔵 TelegramApp: Opening Telegram link...');
      const success = this.openTelegramLink(messageUrl);
      console.log('🔵 TelegramApp: openTelegramLink returned:', success);
      
      if (success) {
        // Provide haptic feedback for successful redirection
        this.hapticFeedback('impact', 'medium');
        console.log('🟢 TelegramApp: Successfully redirected to message:', message.telegram_message_id);
      } else {
        console.error('🔴 TelegramApp: openTelegramLink returned false');
        throw new Error('Failed to open Telegram link');
      }

      return success;
    } catch (error) {
      console.error('🔴 TelegramApp: Exception in redirectToMessage:', error);
      this.showAlert('Unable to open message in Telegram. Please try again.');
      this.hapticFeedback('notification', 'error');
      return false;
    }
  }
}

// Create singleton instance
const telegramApp = new TelegramWebApp();

export default telegramApp;
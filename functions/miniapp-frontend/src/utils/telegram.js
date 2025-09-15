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
      console.log('🔵 TelegramApp: No Telegram WebApp detected - platform unknown');
      return 'unknown';
    }

    const platform = this.tg.platform;
    const version = this.tg.version;
    const colorScheme = this.tg.colorScheme;
    
    console.log('🔵 TelegramApp: Platform detection:', {
      platform,
      version,
      colorScheme,
      userAgent: navigator.userAgent
    });

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
    console.log('🔵 TelegramApp: this.user?.username:', this.user?.username);
    
    if (!messageId) {
      console.error('🔴 TelegramApp: Message ID is required for URL generation');
      return null;
    }

    // Detect platform to choose appropriate URL scheme
    const platform = this.detectPlatform();
    const isNative = this.isNativeClient();
    console.log('🔵 TelegramApp: Detected platform:', platform, 'isNative:', isNative);

    // For native clients (desktop/mobile), try tg:// protocol schemes
    if (isNative && this.user?.username) {
      // Use tg://resolve for native clients with username
      const tgUrl = `tg://resolve?domain=${this.user.username}&post=${messageId}`;
      console.log('🔵 TelegramApp: Generated tg:// URL for native client:', tgUrl);
      return tgUrl;
    }

    // For web clients or fallback, use https://t.me/ URLs
    if (this.user?.username) {
      // Try using username format for saved messages
      const url = `https://t.me/${this.user.username}/${messageId}`;
      console.log('🔵 TelegramApp: Generated https://t.me/ URL for web client:', url);
      return url;
    }
    
    // Last resort fallback: Try the original format
    const targetChatId = chatId || this.user?.id;
    console.log('🔵 TelegramApp: targetChatId resolved to:', targetChatId);
    
    if (!targetChatId) {
      console.error('🔴 TelegramApp: No chat ID available for URL generation');
      console.error('🔴 TelegramApp: chatId provided:', chatId);
      console.error('🔴 TelegramApp: this.user?.id:', this.user?.id);
      return null;
    }

    // Original format as final fallback
    const url = `https://t.me/c/${targetChatId}/${messageId}`;
    console.log('🔵 TelegramApp: Generated fallback URL:', url);
    
    return url;
  }

  /**
   * Redirect to a specific Telegram message with multiple fallback strategies
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

    return this.tryMultipleRedirectionStrategies(message, chatId);
  }

  /**
   * Try multiple redirection strategies for better compatibility
   * @param {Object} message - Message object with telegram_message_id
   * @param {number} chatId - Optional chat ID
   * @returns {boolean} True if any strategy was attempted
   */
  tryMultipleRedirectionStrategies(message, chatId = null) {
    const messageId = message.telegram_message_id;
    const platform = this.detectPlatform();
    const isNative = this.isNativeClient();

    console.log('🔵 TelegramApp: Trying multiple redirection strategies...');
    console.log('🔵 TelegramApp: Platform:', platform, 'IsNative:', isNative);

    // Strategy 1: Platform-specific primary approach
    try {
      const primaryUrl = this.generateMessageUrl(messageId, chatId);
      if (primaryUrl) {
        console.log('🔵 TelegramApp: Strategy 1 - Primary URL:', primaryUrl);
        const success = this.openTelegramLink(primaryUrl);
        if (success) {
          this.hapticFeedback('impact', 'medium');
          console.log('🟢 TelegramApp: Strategy 1 succeeded');
          return true;
        }
      }
    } catch (error) {
      console.warn('🟡 TelegramApp: Strategy 1 failed:', error);
    }

    // Strategy 2: For native clients, try opening just Saved Messages chat
    if (isNative && this.user?.username) {
      try {
        const savedMessagesUrl = `tg://resolve?domain=${this.user.username}`;
        console.log('🔵 TelegramApp: Strategy 2 - Saved Messages chat URL:', savedMessagesUrl);
        const success = this.openTelegramLink(savedMessagesUrl);
        if (success) {
          this.hapticFeedback('impact', 'medium');
          console.log('🟢 TelegramApp: Strategy 2 succeeded - opened Saved Messages chat');
          this.showAlert(`Opened your Saved Messages. Look for message #${messageId}`);
          return true;
        }
      } catch (error) {
        console.warn('🟡 TelegramApp: Strategy 2 failed:', error);
      }
    }

    // Strategy 3: For web clients, try alternative URL format
    if (!isNative && this.user?.username) {
      try {
        const webUrl = `https://t.me/${this.user.username}`;
        console.log('🔵 TelegramApp: Strategy 3 - Web chat URL:', webUrl);
        const success = this.openTelegramLink(webUrl);
        if (success) {
          this.hapticFeedback('impact', 'medium');
          console.log('🟢 TelegramApp: Strategy 3 succeeded - opened personal chat');
          this.showAlert(`Opened your personal chat. Look for message #${messageId}`);
          return true;
        }
      } catch (error) {
        console.warn('🟡 TelegramApp: Strategy 3 failed:', error);
      }
    }

    // Strategy 4: Last resort - inform user about manual navigation
    console.error('🔴 TelegramApp: All redirection strategies failed');
    this.showAlert(`Unable to redirect directly. Please manually go to Saved Messages and look for message ID: ${messageId}`);
    this.hapticFeedback('notification', 'error');
    return false;
  }
}

// Create singleton instance
const telegramApp = new TelegramWebApp();

export default telegramApp;
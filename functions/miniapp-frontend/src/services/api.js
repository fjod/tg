/**
 * API service for communicating with the backend miniapp-api
 */

import telegramApp from '../utils/telegram.js';

class ApiService {
  constructor() {
    // Use API Gateway URL instead of direct function URL for proper HTTP handling
    this.baseURL = 'https://d5di1npf8thkd9m534rv.8wihnuyr.apigw.yandexcloud.net';
    this.timeout = 10000; // 10 seconds timeout
  }

  /**
   * Make authenticated API request
   * @param {string} endpoint - API endpoint (e.g., '/api/user/tags')
   * @param {Object} options - Request options
   * @returns {Promise<Object>} API response
   */
  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    
    // Get authentication data from Telegram WebApp
    const authData = telegramApp.getAuthData();
    if (!authData) {
      throw new Error('Not authenticated - Telegram WebApp data not available');
    }

    const defaultOptions = {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authData}`, // Send initData as Bearer token
      },
      timeout: this.timeout,
    };

    const mergedOptions = {
      ...defaultOptions,
      ...options,
      headers: {
        ...defaultOptions.headers,
        ...options.headers,
      },
    };

    try {
      console.log(`Making API request to: ${url}`, {
        method: mergedOptions.method,
        hasAuth: !!authData,
      });

      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), this.timeout);

      const response = await fetch(url, {
        ...mergedOptions,
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`HTTP ${response.status}: ${errorText || response.statusText}`);
      }

      const data = await response.json();
      console.log(`API response from ${endpoint}:`, data);

      return data;
    } catch (error) {
      console.error(`API request failed for ${endpoint}:`, error);
      
      if (error.name === 'AbortError') {
        throw new Error('Request timeout - please try again');
      }
      
      throw error;
    }
  }

  /**
   * Get user's tags with message counts
   * @returns {Promise<Array>} Array of tag objects
   */
  async getUserTags() {
    try {
      const response = await this.request('/api/user/tags');
      
      if (!response.success) {
        throw new Error(response.error || 'Failed to fetch tags');
      }

      return response.data || [];
    } catch (error) {
      console.error('Failed to get user tags:', error);
      throw new Error(`Failed to load tags: ${error.message}`);
    }
  }

  /**
   * Health check endpoint (for testing)
   * @returns {Promise<boolean>} True if API is accessible
   */
  async healthCheck() {
    try {
      await this.request('/api/health');
      return true;
    } catch (error) {
      console.warn('Health check failed:', error.message);
      return false;
    }
  }

  /**
   * Set custom API base URL (for testing)
   * @param {string} baseURL - New base URL
   */
  setBaseURL(baseURL) {
    this.baseURL = baseURL;
    console.log('API base URL changed to:', baseURL);
  }

  /**
   * Get current API base URL
   * @returns {string} Current base URL
   */
  getBaseURL() {
    return this.baseURL;
  }
}

// Create singleton instance
const apiService = new ApiService();

export default apiService;